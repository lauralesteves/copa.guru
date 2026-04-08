// MongoDB Init Script for Copa.Guru
// Run with: mongosh "your-connection-string" < init.mongodb.js
//
// Creates all collections. Safe to re-run: createCollection is idempotent.

print("======================================");
print("Initializing Copa.Guru database");
print("======================================\n");

const collections = ["users", "matches", "predictions", "leaderboard"];

collections.forEach(name => {
  try {
    db.createCollection(name);
    print("  [create] " + name);
  } catch (e) {
    if (e.codeName === "NamespaceExists") {
      print("  [skip] " + name + " (already exists)");
    } else {
      throw e;
    }
  }
});

print("\n======================================");
print("Verifying collections");
print("======================================\n");

db.getCollectionNames().sort().forEach(name => {
  print("  - " + name);
});

print("\nDone.");
