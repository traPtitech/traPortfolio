package mockdata

import "github.com/gofrs/uuid"

const (
	userID1           uuidStr = "11111111-1111-4111-8111-111111111111"
	userID2           uuidStr = "22222222-2222-4222-8222-222222222222"
	userID3           uuidStr = "33333333-3333-4333-8333-333333333333"
	accountID         uuidStr = "d834e180-2af9-4cfe-838a-8a3930666490"
	contestID         uuidStr = "08eec963-0f29-48d1-929f-004cb67d8ce6"
	contestTeamID     uuidStr = "a9d07124-ffee-412f-adfc-02d3db0b750d"
	groupID           uuidStr = "455938b1-635f-4b43-ae74-66550b04c5d4"
	projectID1        uuidStr = "bf9c1aec-7e3a-4587-8adc-651895aa6ec0"
	projectID2        uuidStr = "bf9c1aec-7e3a-4587-8adc-651895aa6ec1"
	projectID3        uuidStr = "bf9c1aec-7e3a-4587-8adc-651895aa6ec2"
	projectMemberID1  uuidStr = "a211a49c-9b30-48b9-8dbb-c449c99f12c7"
	projectMemberID2  uuidStr = "a211a49c-9b30-48b9-8dbb-c449c99f12c8"
	projectMemberID3  uuidStr = "a211a49c-9b30-48b9-8dbb-c449c99f12c9"
	knoqEventID1      uuidStr = "d1274c6e-15cc-4ca0-b720-1c03ea3a60ec"
	knoqEventID2      uuidStr = "e28ec610-226d-49c5-be7c-86af54f6839d"
	knoqEventGroupID1 uuidStr = "7ecabb2a-8e2c-4ebe-bb0b-13254a6eae05"
	knoqEventGroupID2 uuidStr = "9c592124-52a5-4981-a2c8-1e218c64a8e5"
	knoqEventRoomID1  uuidStr = "68319c0c-be20-45c1-a05d-7651473bd396"
	knoqEventRoomID2  uuidStr = "cbd48b1f-6b20-41c8-b122-a9826bd968ed"

	userName1     = "user1"
	userName2     = "user2"
	userName3     = "lolico"
	userRealname1 = "ユーザー1 ユーザー1"
	userRealname2 = "ユーザー2 ユーザー2"
	userRealname3 = "東 工子"
)

type uuidStr string

func (s uuidStr) uuid() uuid.UUID {
	return uuid.FromStringOrNil(string(s))
}

func UserID1() uuid.UUID {
	return uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111")
}

func UserID2() uuid.UUID {
	return uuid.FromStringOrNil("22222222-2222-4222-8222-222222222222")
}

func UserID3() uuid.UUID {
	return uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333")
}

func AccountID() uuid.UUID {
	return uuid.FromStringOrNil("d834e180-2af9-4cfe-838a-8a3930666490")
}

func ContestID() uuid.UUID {
	return uuid.FromStringOrNil("08eec963-0f29-48d1-929f-004cb67d8ce6")
}

func ContestTeamID() uuid.UUID {
	return uuid.FromStringOrNil("a9d07124-ffee-412f-adfc-02d3db0b750d")
}

func GroupID() uuid.UUID {
	return uuid.FromStringOrNil("455938b1-635f-4b43-ae74-66550b04c5d4")
}

func ProjectID1() uuid.UUID {
	return uuid.FromStringOrNil("bf9c1aec-7e3a-4587-8adc-651895aa6ec0")
}

func ProjectID2() uuid.UUID {
	return uuid.FromStringOrNil("bf9c1aec-7e3a-4587-8adc-651895aa6ec1")
}

func ProjectID3() uuid.UUID {
	return uuid.FromStringOrNil("bf9c1aec-7e3a-4587-8adc-651895aa6ec2")
}

func ProjectMemberID1() uuid.UUID {
	return uuid.FromStringOrNil("a211a49c-9b30-48b9-8dbb-c449c99f12c7")
}

func ProjectMemberID2() uuid.UUID {
	return uuid.FromStringOrNil("a211a49c-9b30-48b9-8dbb-c449c99f12c8")
}

func ProjectMemberID3() uuid.UUID {
	return uuid.FromStringOrNil("a211a49c-9b30-48b9-8dbb-c449c99f12c9")
}

func KnoqEventID1() uuid.UUID {
	return uuid.FromStringOrNil("d1274c6e-15cc-4ca0-b720-1c03ea3a60ec")
}

func KnoqEventID2() uuid.UUID {
	return uuid.FromStringOrNil("e28ec610-226d-49c5-be7c-86af54f6839d")
}

func KnoqEventGroupID1() uuid.UUID {
	return uuid.FromStringOrNil("7ecabb2a-8e2c-4ebe-bb0b-13254a6eae05")
}

func KnoqEventGroupID2() uuid.UUID {
	return uuid.FromStringOrNil("9c592124-52a5-4981-a2c8-1e218c64a8e5")
}

func KnoqEventRoomID1() uuid.UUID {
	return uuid.FromStringOrNil("68319c0c-be20-45c1-a05d-7651473bd396")
}

func KnoqEventRoomID2() uuid.UUID {
	return uuid.FromStringOrNil("cbd48b1f-6b20-41c8-b122-a9826bd968ed")
}
