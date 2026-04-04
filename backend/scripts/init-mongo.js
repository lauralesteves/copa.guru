db = db.getSiblingDB("copa-guru");

db.createCollection("users");
db.createCollection("matches");
db.createCollection("predictions");
db.createCollection("leaderboard");

db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ googleId: 1 }, { unique: true });
db.matches.createIndex({ date: 1 });
db.predictions.createIndex({ user_id: 1, match_id: 1 }, { unique: true });
db.leaderboard.createIndex({ points: -1 });

print("copa-guru database initialized");
