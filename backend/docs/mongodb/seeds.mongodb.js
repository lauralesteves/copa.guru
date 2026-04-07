// MongoDB Seed Script for Copa.Guru (local development)
// Run with: mongosh "mongodb://localhost:27017/copa-guru" < seed.mongodb.js
//
// Creates 5 test users, 6 matches (4 finished, 1 live, 1 scheduled),
// predictions for each user, and a leaderboard.
// Safe to re-run: drops collections before inserting.

print("======================================");
print("Seeding Copa.Guru local database");
print("======================================\n");

const now = new Date();
const passwordHash = "$2y$10$YMN//pdTzGKd4ICwMMGHqeS48LMO426coFGFkj6N/Q.QyzfsPofCC"; // 123password

// ============================================
// USERS
// ============================================

print("Seeding users...");
db.users.drop();

const users = db.users.insertMany([
  {
    name: "Test User 1", email: "test1@copa.guru", picture: "",
    strategy: "email", auth: { password: passwordHash, lastLoginAt: now },
    createdAt: now, updatedAt: now
  },
  {
    name: "Test User 2", email: "test2@copa.guru", picture: "",
    strategy: "email", auth: { password: passwordHash, lastLoginAt: now },
    createdAt: now, updatedAt: now
  },
  {
    name: "Test User 3", email: "test3@copa.guru", picture: "",
    strategy: "email", auth: { password: passwordHash, lastLoginAt: now },
    createdAt: now, updatedAt: now
  },
  {
    name: "Test User 4", email: "test4@copa.guru", picture: "",
    strategy: "email", auth: { password: passwordHash, lastLoginAt: now },
    createdAt: now, updatedAt: now
  },
  {
    name: "Test User 5", email: "test5@copa.guru", picture: "",
    strategy: "email", auth: { password: passwordHash, lastLoginAt: now },
    createdAt: now, updatedAt: now
  }
]);

const userIds = Object.values(users.insertedIds);
print("  " + userIds.length + " users created");

// ============================================
// MATCHES
// ============================================

print("Seeding matches...");
db.matches.drop();

const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000);
const twoDaysAgo = new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000);
const threeDaysAgo = new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000);
const fourDaysAgo = new Date(now.getTime() - 4 * 24 * 60 * 60 * 1000);
const tomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000);

const matches = db.matches.insertMany([
  {
    matchNumber: 1, team1: "Brasil", team2: "Sérvia",
    goals1: 2, goals2: 0, date: fourDaysAgo,
    stadium: "Estádio Lusail", city: "Lusail", group: "G", stage: "GROUP_STAGE",
    status: "FINISHED", updatedAt: now
  },
  {
    matchNumber: 2, team1: "Argentina", team2: "Arábia Saudita",
    goals1: 1, goals2: 2, date: threeDaysAgo,
    stadium: "Estádio Lusail", city: "Lusail", group: "C", stage: "GROUP_STAGE",
    status: "FINISHED", updatedAt: now
  },
  {
    matchNumber: 3, team1: "Alemanha", team2: "Japão",
    goals1: 1, goals2: 2, date: twoDaysAgo,
    stadium: "Khalifa International", city: "Al Rayyan", group: "E", stage: "GROUP_STAGE",
    status: "FINISHED", updatedAt: now
  },
  {
    matchNumber: 4, team1: "Espanha", team2: "Costa Rica",
    goals1: 7, goals2: 0, date: yesterday,
    stadium: "Al Thumama", city: "Doha", group: "E", stage: "GROUP_STAGE",
    status: "FINISHED", updatedAt: now
  },
  {
    matchNumber: 5, team1: "França", team2: "Austrália",
    goals1: 2, goals2: 1, date: now,
    stadium: "Al Janoub", city: "Al Wakrah", group: "D", stage: "GROUP_STAGE",
    status: "LIVE", updatedAt: now
  },
  {
    matchNumber: 6, team1: "Portugal", team2: "Gana",
    goals1: null, goals2: null, date: tomorrow,
    stadium: "Stadium 974", city: "Doha", group: "H", stage: "GROUP_STAGE",
    status: "SCHEDULED", updatedAt: now
  }
]);

const matchIds = Object.values(matches.insertedIds);
print("  " + matchIds.length + " matches created");

// ============================================
// PREDICTIONS (only for FINISHED matches)
// ============================================

print("Seeding predictions...");
db.predictions.drop();

// Match 1: Brasil 2x0 Sérvia
// Match 2: Argentina 1x2 Arábia Saudita
// Match 3: Alemanha 1x2 Japão
// Match 4: Espanha 7x0 Costa Rica

const predictions = [
  // --- User 1: 1 exact + 2 correct + 1 wrong = 20 pts ---
  { userId: userIds[0], matchId: matchIds[0], matchNumber: 1, goals1: 2, goals2: 0, points: 10, createdAt: now, updatedAt: now }, // exact
  { userId: userIds[0], matchId: matchIds[1], matchNumber: 2, goals1: 0, goals2: 1, points: 5,  createdAt: now, updatedAt: now }, // correct result
  { userId: userIds[0], matchId: matchIds[2], matchNumber: 3, goals1: 0, goals2: 1, points: 5,  createdAt: now, updatedAt: now }, // correct result
  { userId: userIds[0], matchId: matchIds[3], matchNumber: 4, goals1: 2, goals2: 1, points: 0,  createdAt: now, updatedAt: now }, // wrong

  // --- User 2: 2 exact + 1 correct + 1 wrong = 25 pts ---
  { userId: userIds[1], matchId: matchIds[0], matchNumber: 1, goals1: 2, goals2: 0, points: 10, createdAt: now, updatedAt: now }, // exact
  { userId: userIds[1], matchId: matchIds[1], matchNumber: 2, goals1: 1, goals2: 2, points: 10, createdAt: now, updatedAt: now }, // exact
  { userId: userIds[1], matchId: matchIds[2], matchNumber: 3, goals1: 2, goals2: 3, points: 5,  createdAt: now, updatedAt: now }, // correct result
  { userId: userIds[1], matchId: matchIds[3], matchNumber: 4, goals1: 1, goals2: 1, points: 0,  createdAt: now, updatedAt: now }, // wrong

  // --- User 3: 0 exact + 3 correct + 1 wrong = 15 pts ---
  { userId: userIds[2], matchId: matchIds[0], matchNumber: 1, goals1: 1, goals2: 0, points: 5,  createdAt: now, updatedAt: now }, // correct result
  { userId: userIds[2], matchId: matchIds[1], matchNumber: 2, goals1: 0, goals2: 3, points: 5,  createdAt: now, updatedAt: now }, // correct result
  { userId: userIds[2], matchId: matchIds[2], matchNumber: 3, goals1: 0, goals2: 3, points: 5,  createdAt: now, updatedAt: now }, // correct result
  { userId: userIds[2], matchId: matchIds[3], matchNumber: 4, goals1: 0, goals2: 1, points: 0,  createdAt: now, updatedAt: now }, // wrong

  // --- User 4: 1 exact + 0 correct + 3 wrong = 10 pts ---
  { userId: userIds[3], matchId: matchIds[0], matchNumber: 1, goals1: 0, goals2: 2, points: 0,  createdAt: now, updatedAt: now }, // wrong
  { userId: userIds[3], matchId: matchIds[1], matchNumber: 2, goals1: 2, goals2: 0, points: 0,  createdAt: now, updatedAt: now }, // wrong
  { userId: userIds[3], matchId: matchIds[2], matchNumber: 3, goals1: 1, goals2: 2, points: 10, createdAt: now, updatedAt: now }, // exact
  { userId: userIds[3], matchId: matchIds[3], matchNumber: 4, goals1: 0, goals2: 1, points: 0,  createdAt: now, updatedAt: now }, // wrong

  // --- User 5: 0 exact + 1 correct + 3 wrong = 5 pts ---
  { userId: userIds[4], matchId: matchIds[0], matchNumber: 1, goals1: 0, goals2: 1, points: 0,  createdAt: now, updatedAt: now }, // wrong
  { userId: userIds[4], matchId: matchIds[1], matchNumber: 2, goals1: 3, goals2: 0, points: 0,  createdAt: now, updatedAt: now }, // wrong
  { userId: userIds[4], matchId: matchIds[2], matchNumber: 3, goals1: 2, goals2: 0, points: 0,  createdAt: now, updatedAt: now }, // wrong
  { userId: userIds[4], matchId: matchIds[3], matchNumber: 4, goals1: 5, goals2: 0, points: 5,  createdAt: now, updatedAt: now }, // correct result

  // --- Match 6 (SCHEDULED): User 1 and User 2 have future predictions (no points yet) ---
  { userId: userIds[0], matchId: matchIds[5], matchNumber: 6, goals1: 3, goals2: 1, points: null, createdAt: now, updatedAt: now },
  { userId: userIds[1], matchId: matchIds[5], matchNumber: 6, goals1: 2, goals2: 0, points: null, createdAt: now, updatedAt: now }
];

db.predictions.insertMany(predictions);
print("  " + predictions.length + " predictions created");

// ============================================
// LEADERBOARD
// ============================================

print("Seeding leaderboard...");
db.leaderboard.drop();

db.leaderboard.insertMany([
  {
    userId: userIds[1], userName: "Test User 2", userPicture: "",
    totalPoints: 25, exactScores: 2, correctResults: 1, totalPredictions: 4, rank: 1,
    updatedAt: now
  },
  {
    userId: userIds[0], userName: "Test User 1", userPicture: "",
    totalPoints: 20, exactScores: 1, correctResults: 2, totalPredictions: 4, rank: 2,
    updatedAt: now
  },
  {
    userId: userIds[2], userName: "Test User 3", userPicture: "",
    totalPoints: 15, exactScores: 0, correctResults: 3, totalPredictions: 4, rank: 3,
    updatedAt: now
  },
  {
    userId: userIds[3], userName: "Test User 4", userPicture: "",
    totalPoints: 10, exactScores: 1, correctResults: 0, totalPredictions: 4, rank: 4,
    updatedAt: now
  },
  {
    userId: userIds[4], userName: "Test User 5", userPicture: "",
    totalPoints: 5, exactScores: 0, correctResults: 1, totalPredictions: 4, rank: 5,
    updatedAt: now
  }
]);

print("  5 leaderboard entries created");

// ============================================
// INDEXES
// ============================================

print("\nCreating indexes...");
load("docs/mongodb/indexes.mongodb.js");

print("\n======================================");
print("Seed complete!");
print("======================================");
print("\nRanking:");
print("  1. Test User 2 — 25 pts (2 exact, 1 correct)");
print("  2. Test User 1 — 20 pts (1 exact, 2 correct)");
print("  3. Test User 3 — 15 pts (0 exact, 3 correct)");
print("  4. Test User 4 — 10 pts (1 exact, 0 correct)");
print("  5. Test User 5 —  5 pts (0 exact, 1 correct)");
print("\nLogin: test[1-5]@copa.guru / 123password");
