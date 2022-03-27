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
