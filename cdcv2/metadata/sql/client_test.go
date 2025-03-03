// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pingcap/tiflow/pkg/errors"
	"github.com/stretchr/testify/require"
)

// Note that updateAt is not included in the test because it is automatically updated by gorm.
// TODO(CharlesCheung): add test to verify the correctness of updateAt.

func runMockExecTest(t *testing.T, mock sqlmock.Sqlmock, expectedSQL string, args []driver.Value, fn func() error) {
	testErr := errors.New("test error")

	// Test normal execution
	mock.ExpectExec(expectedSQL).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))
	err := fn()
	require.NoError(t, err)

	// Test rows affected not match
	mock.ExpectExec(expectedSQL).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 0))
	err = fn()
	require.ErrorIs(t, err, errors.ErrMetaRowsAffectedNotMatch)

	// Test op failed
	mock.ExpectExec(expectedSQL).WithArgs(args...).WillReturnError(testErr)
	err = fn()
	require.ErrorIs(t, err, errors.ErrMetaOpFailed)
	require.ErrorContains(t, err, testErr.Error())
}

func TestUpstreamClientExecSQL(t *testing.T) {
	t.Parallel()

	backendDB, db, mock := newMockDB(t)
	defer backendDB.Close()
	cient := NewORMClient("test-upstream-client", db)

	up := &UpstreamDO{
		ID:        1,
		Endpoints: "endpoints",
		Config: &Credential{
			CAPath: "ca-path",
		},
		Version: 1,
	}
	config, err := up.Config.Value()
	require.NoError(t, err)

	// Test createUpstream
	runMockExecTest(
		t, mock,
		"INSERT INTO `upstream` (`endpoints`,`config`,`version`,`update_at`,`id`) VALUES (?,?,?,?,?)",
		[]driver.Value{up.Endpoints, up.Config, up.Version, sqlmock.AnyArg(), up.ID},
		func() error {
			return cient.createUpstream(db, up)
		},
	)

	// Test deleteUpstream
	runMockExecTest(
		t, mock,
		// TODO(CharlesCheung): delete statement should be optimized, such as remove duplicated fields.
		// Note: should the version be checked?
		"DELETE FROM `upstream` WHERE `upstream`.`id` = ?",
		[]driver.Value{up.ID},
		func() error {
			return cient.deleteUpstream(db, up)
		},
	)

	// Test updateUpstream
	runMockExecTest(
		t, mock,
		"UPDATE `upstream` SET `endpoints`=?,`config`=?,`version`=?,`update_at`=? WHERE id = ? and version = ?",
		[]driver.Value{up.Endpoints, config, up.Version + 1, sqlmock.AnyArg(), up.ID, up.Version},
		func() error {
			return cient.updateUpstream(db, up)
		},
	)

	// Test updateUpstream with nil config
	up.Config = nil
	runMockExecTest(
		t, mock,
		"UPDATE `upstream` SET `endpoints`=?,`version`=?,`update_at`=? WHERE id = ? and version = ?",
		[]driver.Value{up.Endpoints, up.Version + 1, sqlmock.AnyArg(), up.ID, up.Version},
		func() error {
			return cient.updateUpstream(db, up)
		},
	)
}

func TestChangefeedInfoClientExecSQL(t *testing.T) {
	t.Parallel()

	backendDB, db, mock := newMockDB(t)
	defer backendDB.Close()
	cient := NewORMClient("test-changefeed-info-client", db)

	info := &ChangefeedInfoDO{
		UUID:       1,
		Namespace:  "namespace",
		ID:         "id",
		RemovedAt:  nil,
		UpstreamID: 1,
		SinkURI:    "sinkURI",
		StartTs:    1,
		TargetTs:   1,
		Config:     &ReplicaConfig{},
		Version:    1,
	}
	configValue, err := info.Config.Value()
	require.NoError(t, err)

	// Test createChangefeedInfo
	runMockExecTest(
		t, mock,
		"INSERT INTO `changefeed_info` ("+
			"`namespace`,`id`,`removed_at`,`upstream_id`,"+
			"`sink_uri`,`start_ts`,`target_ts`,`config`,"+
			"`version`,`update_at`,`uuid`) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		[]driver.Value{
			info.Namespace, info.ID, info.RemovedAt, info.UpstreamID,
			info.SinkURI, info.StartTs, info.TargetTs, configValue,
			info.Version, sqlmock.AnyArg(), info.UUID,
		},
		func() error {
			return cient.createChangefeedInfo(db, info)
		},
	)

	// Test deleteChangefeedInfo
	runMockExecTest(
		t, mock,
		"DELETE FROM `changefeed_info` WHERE `changefeed_info`.`uuid` = ?",
		[]driver.Value{info.UUID},
		func() error {
			return cient.deleteChangefeedInfo(db, info)
		},
	)

	// Test updateChangefeedInfo
	runMockExecTest(
		t, mock,
		"UPDATE `changefeed_info` SET `sink_uri`=?,`start_ts`=?,`target_ts`=?,`config`=?,`version`=?,`update_at`=? WHERE uuid = ? and version = ?",
		[]driver.Value{info.SinkURI, info.StartTs, info.TargetTs, configValue, info.Version + 1, sqlmock.AnyArg(), info.UUID, info.Version},
		func() error {
			return cient.updateChangefeedInfo(db, info)
		},
	)

	// Test updateChangefeedInfo with nil config
	info.Config = nil
	runMockExecTest(
		t, mock,
		"UPDATE `changefeed_info` SET `sink_uri`=?,`start_ts`=?,`target_ts`=?,`version`=?,`update_at`=? WHERE uuid = ? and version = ?",
		[]driver.Value{info.SinkURI, info.StartTs, info.TargetTs, info.Version + 1, sqlmock.AnyArg(), info.UUID, info.Version},
		func() error {
			return cient.updateChangefeedInfo(db, info)
		},
	)
}

func TestChangefeedStateClientExecSQL(t *testing.T) {
	t.Parallel()

	backendDB, db, mock := newMockDB(t)
	defer backendDB.Close()
	cient := NewORMClient("test-changefeed-state-client", db)

	state := &ChangefeedStateDO{
		ChangefeedUUID: 1,
		State:          "state",
		// Note that warning and error could be nil.
		Warning: nil,
		Error: &RunningError{
			Time: time.Now(),
			Addr: "addr",
			Code: "code",
		},
		Version: 1,
	}

	errVal, err := state.Error.Value()
	require.NoError(t, err)

	// Test createChangefeedState
	runMockExecTest(
		t, mock,
		"INSERT INTO `changefeed_state` (`state`,`warning`,`error`,`version`,`update_at`,`changefeed_uuid`) VALUES (?,?,?,?,?,?)",
		[]driver.Value{state.State, state.Warning, errVal, state.Version, sqlmock.AnyArg(), state.ChangefeedUUID},
		func() error {
			return cient.createChangefeedState(db, state)
		},
	)

	// Test deleteChangefeedState
	runMockExecTest(
		t, mock,
		"DELETE FROM `changefeed_state` WHERE `changefeed_state`.`changefeed_uuid` = ?",
		[]driver.Value{state.ChangefeedUUID},
		func() error {
			return cient.deleteChangefeedState(db, state)
		},
	)

	// Test updateChangefeedState
	// Note that a nil error or warning will also be updated.
	runMockExecTest(
		t, mock,
		"UPDATE `changefeed_state` SET `state`=?,`warning`=?,`error`=?,`version`=?,`update_at`=? WHERE changefeed_uuid = ? and version = ?",
		[]driver.Value{state.State, state.Warning, errVal, state.Version + 1, sqlmock.AnyArg(), state.ChangefeedUUID, state.Version},
		func() error {
			return cient.updateChangefeedState(db, state)
		},
	)
}

func TestScheduleClientExecSQL(t *testing.T) {
	t.Parallel()

	backendDB, db, mock := newMockDB(t)
	defer backendDB.Close()
	cient := NewORMClient("test-schedule-client", db)

	ownerCapture := "test-owner"
	schedule := &ScheduleDO{
		ChangefeedUUID: 1,
		Owner:          &ownerCapture,
		OwnerState:     "ownerState",
		Processors:     nil,
		TaskPosition: ChangefeedProgress{
			CheckpointTs: 1,
		},
		Version: 1,
	}

	// Test createSchedule
	runMockExecTest(
		t, mock,
		"INSERT INTO `schedule` (`owner`,`owner_state`,`processors`,`task_position`,`version`,`update_at`,`changefeed_uuid`) VALUES (?,?,?,?,?,?,?)",
		[]driver.Value{schedule.Owner, schedule.OwnerState, schedule.Processors, schedule.TaskPosition, schedule.Version, sqlmock.AnyArg(), schedule.ChangefeedUUID},
		func() error {
			return cient.createSchedule(db, schedule)
		},
	)

	// Test deleteSchedule
	runMockExecTest(
		t, mock,
		"DELETE FROM `schedule` WHERE `schedule`.`changefeed_uuid` = ?",
		[]driver.Value{schedule.ChangefeedUUID},
		func() error {
			return cient.deleteSchedule(db, schedule)
		},
	)

	// Test updateSchedule with non-empty task position.
	runMockExecTest(
		t, mock,
		"UPDATE `schedule` SET `owner`=?,`owner_state`=?,`task_position`=?,`version`=?,`update_at`=? WHERE changefeed_uuid = ? and version = ?",
		[]driver.Value{schedule.Owner, schedule.OwnerState, schedule.TaskPosition, schedule.Version + 1, sqlmock.AnyArg(), schedule.ChangefeedUUID, schedule.Version},
		func() error {
			return cient.updateSchedule(db, schedule)
		},
	)

	// Test updateSchedule with empty task position.
	schedule.TaskPosition = ChangefeedProgress{}
	runMockExecTest(
		t, mock,
		"UPDATE `schedule` SET `owner`=?,`owner_state`=?,`version`=?,`update_at`=? WHERE changefeed_uuid = ? and version = ?",
		[]driver.Value{schedule.Owner, schedule.OwnerState, schedule.Version + 1, sqlmock.AnyArg(), schedule.ChangefeedUUID, schedule.Version},
		func() error {
			return cient.updateSchedule(db, schedule)
		},
	)
}

func TestProgressClientExecSQL(t *testing.T) {
	t.Parallel()

	backendDB, db, mock := newMockDB(t)
	defer backendDB.Close()
	cient := NewORMClient("test-progress-client", db)

	progress := &ProgressDO{
		CaptureID: "test-captureID",
		Progress:  nil,
		Version:   1,
	}

	// Test createProgress
	runMockExecTest(
		t, mock,
		"INSERT INTO `progress` (`capture_id`,`progress`,`version`,`update_at`) VALUES (?,?,?,?)",
		[]driver.Value{progress.CaptureID, progress.Progress, progress.Version, sqlmock.AnyArg()},
		func() error {
			return cient.createProgress(db, progress)
		},
	)

	// Test deleteProgress
	runMockExecTest(
		t, mock,
		"DELETE FROM `progress` WHERE `progress`.`capture_id` = ?",
		[]driver.Value{progress.CaptureID},
		func() error {
			return cient.deleteProgress(db, progress)
		},
	)

	// Test updateProgress
	progress.Progress = &CaptureProgress{}
	runMockExecTest(
		t, mock,
		"UPDATE `progress` SET `progress`=?,`version`=?,`update_at`=? WHERE capture_id = ? and version = ?",
		[]driver.Value{progress.Progress, progress.Version + 1, sqlmock.AnyArg(), progress.CaptureID, progress.Version},
		func() error {
			return cient.updateProgress(db, progress)
		},
	)
}
