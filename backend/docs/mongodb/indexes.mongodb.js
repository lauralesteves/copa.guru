// MongoDB Index Creation Script for Copa.Guru
// Run with: mongosh "your-connection-string" --eval "use copa-guru" < indexes.mongodb.js
//
// Safe to re-run: checks if each index exists before creating.

function ensureIndex(collection, keys, options) {
  const existing = collection.getIndexes().map(i => i.name);
  if (existing.includes(options.name)) {
    print("  [skip] " + collection.getFullName() + "." + options.name + " (already exists)");
    return;
  }
  collection.createIndex(keys, options);
  print("  [create] " + collection.getFullName() + "." + options.name);
}

print("======================================");
print("Creating indexes for Copa.Guru");
print("======================================\n");

// ============================================
// USERS
// ============================================

print("Users collection:\n");

// GetByEmail / Upsert: lookup by email
ensureIndex(db.users, { "email": 1 }, {
  unique: true,
  name: "idx_email_unique"
});

// GetByGoogleID: lookup by google auth id
ensureIndex(db.users, { "auth.googleId": 1 }, {
  name: "idx_auth_google_id",
  sparse: true
});

// GetByRefreshToken: lookup for token refresh flow
ensureIndex(db.users, { "auth.refreshToken": 1 }, {
  name: "idx_auth_refresh_token",
  sparse: true
});

// ============================================
// MATCHES
// ============================================

print("\nMatches collection:\n");

// FindByMatchNumber: unique match identifier
ensureIndex(db.matches, { "matchNumber": 1 }, {
  unique: true,
  name: "idx_match_number_unique"
});

// FindByExternalID: sync with football-data.org
ensureIndex(db.matches, { "externalId": 1 }, {
  unique: true,
  name: "idx_external_id_unique",
  sparse: true
});

// FindByStatus / sync_results worker: filter matches by status
ensureIndex(db.matches, { "status": 1 }, {
  name: "idx_status"
});

// FindScheduledBefore: find matches that should have started
ensureIndex(db.matches, { "date": 1, "status": 1 }, {
  name: "idx_date_status"
});

// List filtering by group phase
ensureIndex(db.matches, { "group": 1 }, {
  name: "idx_group",
  sparse: true
});

// List filtering by stage (GROUP_STAGE, ROUND_OF_16, etc.)
ensureIndex(db.matches, { "stage": 1 }, {
  name: "idx_stage"
});

// ============================================
// PREDICTIONS
// ============================================

print("\nPredictions collection:\n");

// Unique constraint: one prediction per user per match
ensureIndex(db.predictions, { "userId": 1, "matchId": 1 }, {
  unique: true,
  name: "idx_user_match_unique"
});

// FindByMatch / ScoreMatch: all predictions for a given match
ensureIndex(db.predictions, { "matchId": 1 }, {
  name: "idx_match_id"
});

// FindByUser / ListMyPredictions: all predictions for a given user
ensureIndex(db.predictions, { "userId": 1 }, {
  name: "idx_user_id"
});

// ============================================
// LEADERBOARD
// ============================================

print("\nLeaderboard collection:\n");

// GetByUserID: lookup individual user position
ensureIndex(db.leaderboard, { "userId": 1 }, {
  unique: true,
  name: "idx_user_id_unique"
});

// Tiebreaker sorting: points > exact scores
ensureIndex(db.leaderboard, { "totalPoints": -1, "exactScores": -1 }, {
  name: "idx_points_tiebreaker"
});

// List: ordered by rank
ensureIndex(db.leaderboard, { "rank": 1 }, {
  name: "idx_rank"
});

// ============================================
// VERIFICATION
// ============================================

print("\n======================================");
print("Verifying indexes");
print("======================================\n");

["users", "matches", "predictions", "leaderboard"].forEach(name => {
  print(name + ":");
  db[name].getIndexes().forEach(idx => {
    print("  - " + idx.name + (idx.unique ? " (unique)" : "") + (idx.sparse ? " (sparse)" : ""));
  });
  print("");
});

print("Done.");
