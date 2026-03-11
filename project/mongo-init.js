db = db.getSiblingDB('app');

db.createUser({
    user: 'user',
    pwd: 'password',
    roles: [{ role: 'readWrite', db: 'app' }],
});

// Users
db.createCollection('users');
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "username": 1 }, { unique: true });

// Sessions
db.createCollection('sessions');
db.sessions.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 });
db.sessions.createIndex({ "refresh_token_hash": 1 }, { unique: true });
db.sessions.createIndex({ "device_session_id": 1 }, { unique: true });
db.sessions.createIndex({ "user_id": 1 });

// chats
db.chats.createIndex({ "members": 1 })
db.chats.createIndex({ "type": 1 })

// messages
db.messages.createIndex({ "chat_id": 1, "created_at": -1 })
db.messages.createIndex({ "user_id": 1 })

// files
db.files.createIndex({ "chat_id": 1 })