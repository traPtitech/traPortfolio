package mockdata

import "github.com/gofrs/uuid"

const (
	userName1     = "user1"
	userName2     = "user2"
	userName3     = "lolico"
	userRealname1 = "ユーザー1 ユーザー1"
	userRealname2 = "ユーザー2 ユーザー2"
	userRealname3 = "東 工子"
)

var (
	userID1           = uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")
	userID2           = uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222")
	userID3           = uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333")
	accountID         = uuid.FromStringOrNil("d834e180-2af9-4cfe-838a-8a3930666490")
	contestID         = uuid.FromStringOrNil("08eec963-0f29-48d1-929f-004cb67d8ce6")
	contestTeamID     = uuid.FromStringOrNil("a9d07124-ffee-412f-adfc-02d3db0b750d")
	eventID           = uuid.FromStringOrNil("e32a0431-aa0e-4825-98e6-479912275bbd")
	groupID           = uuid.FromStringOrNil("455938b1-635f-4b43-ae74-66550b04c5d4")
	projectID         = uuid.FromStringOrNil("bf9c1aec-7e3a-4587-8adc-651895aa6ec0")
	projectMemberID   = uuid.FromStringOrNil("a211a49c-9b30-48b9-8dbb-c449c99f12c7")
	knoqEventID1      = uuid.FromStringOrNil("d1274c6e-15cc-4ca0-b720-1c03ea3a60ec")
	knoqEventID2      = uuid.FromStringOrNil("e28ec610-226d-49c5-be7c-86af54f6839d")
	knoqEventGroupID1 = uuid.FromStringOrNil("7ecabb2a-8e2c-4ebe-bb0b-13254a6eae05")
	knoqEventGroupID2 = uuid.FromStringOrNil("9c592124-52a5-4981-a2c8-1e218c64a8e5")
	knoqEventRoomID1  = uuid.FromStringOrNil("68319c0c-be20-45c1-a05d-7651473bd396")
	knoqEventRoomID2  = uuid.FromStringOrNil("cbd48b1f-6b20-41c8-b122-a9826bd968ed")
)
