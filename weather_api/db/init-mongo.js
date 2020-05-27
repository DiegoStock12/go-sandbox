db.createUser(
    {
        user: "diego",
        pwd: "passwd",
        roles: [
            {
                role: "readWrite",
                db: "test"
            }
        ]
    }
)