## Create mongo user
```
db.createUser({
 user: "rakutanBot",
 pwd: "rakutanRW",
 roles: [
  {
   role: "readWrite",
   db: "common"
  },
  
 ]
})
```

## Aggregate
```
db.counter.aggregate([{$group: {"_id": "1",rakutan:{$sum: "$normalomikuji"}, onitan: {$sum: "$oniomikuji"}}}]);
{ _id: '1', rakutan: 92927, onitan: 17471 }
db.usertable.aggregate([{$group: {"_id": "1", count: {$sum: "$count"}}}]);
{ _id: '1', count: 946824 }
```
