ColName: db.ColName_Club,
    Keys: bson.D{
    	{Key: "club_id", Value: 1},
    },
    Keys: bson.D{
    	{Key: "invitation_code", Value: 1},
    },


ColName: db.ColName_Club_User_Info,
    Keys: bson.D{
    	{Key: "account", Value: 1},
    },
    Keys: bson.D{
    	{Key: "user_id", Value: 1},
    },
    Keys: bson.D{
    	{Key: "token", Value: 1},
    },


ColName: db.ColName_Table,
    Keys: bson.D{
    	{Key: "table_id", Value: 1},
    },

ColName: db.ColName_GameWallet,
    Keys: bson.D{
    	{Key: "user_id", Value: 1},
    	{Key: "table_id", Value: 1},
    },




