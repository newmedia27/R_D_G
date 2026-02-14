db = db.getSiblingDB('app');

db.createUser({
    user: 'user',
    pwd: 'password',
    roles: [
        {
            role: 'readWrite',
            db: 'app',
        },
    ],
});

// db.createCollection('users');
//
// db.users.createIndex({ createdAt: 1 });
