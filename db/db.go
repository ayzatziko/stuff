package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayzatziko/stuff/xerrors"
)

func Open(ctx context.Context, dsn string) (_ *sql.DB, err error) {
	defer xerrors.Wrap(&err, "db: Open(%s)", dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func LockObjectMigrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `create table lock_object_example(
	id bigserial,
	content text,
	version int
);`)
	if err != nil {
		return fmt.Errorf("db: LockObjectMigrate: %w", err)
	}

	return nil
}

type LockObject struct {
	ID      int
	Content string
	Version int
}

func OptimisticUpdateOfLockObject(ctx context.Context, db *sql.DB, id, cliVer int, f func(LockObject) (LockObject, error)) (err error) {
	defer xerrors.Wrap(&err, "db: OptimisticUpdateOfLockObject(%d)", id)

	row := db.QueryRowContext(ctx, `select content, version from lock_object_example where id = $1`, id)
	oldObj := LockObject{ID: id}
	if err := row.Scan(&oldObj.Content, &oldObj.Version); err != nil {
		return err
	}

	newObject, err := f(oldObj)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		}

		tx.Rollback()
	}()

	row = tx.QueryRowContext(ctx, `select version from lock_object_example where id = $1`, id)
	var curVer int
	if err := row.Scan(&curVer); err != nil {
		return err
	}

	if curVer != cliVer {
		return errors.New("conflict")
	}

	_, err = tx.ExecContext(ctx,
		`update lock_object_example set content = $1, version = $2 where id = $3`,
		newObject.Content,
		curVer+1,
		id,
	)

	return err
}
