# **GameAPI:**  
**On all game api paths:**    
**Header must contain [Authorization: Bearer <accessToken>]**  

- "POST /api/game/create/solo" - creates new solo game with computer  
Response Body:  
```json
{
  "gameuuid": "GameUUID, but it's always == YourUUID, when you play solo",
  "gamefield": {
    "field": [
      [0, 0, 0],
      [0, 0, 0],
      [0, 0, 0]
    ]
  },
  "state": 0 // gameStatusSolo
}
```

- "POST /api/game/create/pair" - creates new pair game with other player  
Response Body:  
```json
{
  "gameuuid": "GameUUID",
  "fPlayerUuid": "UUID of player who created this pair game",
  "sPlayerUuid": "Nil when created, valid when someone joined",
  "gamefield": {
    "field": [
      [0, 0, 0],
      [0, 0, 0],
      [0, 0, 0]
    ]
  },
  "state": 0 // gameStatusPair
}  
```

- "GET /api/game/pair" - get list of all available (joinable) pair games  
Response Body:
```json
{
  "games": ["GameUUID1", "...", ...]
}
```

- "POST /api/game/{GameUUID}/join" - join to a pair game  
Response Body:  
```json
{
  "gameuuid": "GameUUID",
  "fPlayerUuid": "UUID of player who created this pair game",
  "sPlayerUuid": "Your UUID, if you joined this game",
  "gamefield": {
    "field": [
      [0, 0, 0],
      [0, 0, 0],
      [0, 0, 0]
    ]
  },
  "state": 0 // gameStatusPair
} 
```

- "POST /api/game/{GameUUID}" - make a move (pair or solo game)  
Request Body:
```json
{
  "field": [
    [...],
    [...],
    [...]
  ]
}  
```
Response Body Solo:  
```json
{
  "gamefield": {
    "field": [
      [...],[...],[...]
    ]
  },
  "state": 0 // StatusSolo
}  
```
Response Body Pair:  
```json
{
  "gamefield": {
    "field": [
      [...],[...],[...]
    ]
  },
  "state": 0 // StatusPair
}  
```

- "GET /api/game/{GameUUID}/shortinfo" - get pair game shortinfo  
Response Body:  
```json
{
  "gamefield": {
    "field": [
      [...],[...],[...]
    ]
  },
  "state": 0 // gameStatus
}  
```

- "GET /api/game/{GameUUID}/fullinfo" - get pair game fullinfo  
Response Body:  
```json
{
  "gameuuid": "GameUUID",
  "fPlayerUuid": "UUID of player who created this pair game",
  "sPlayerUuid": "UUID of player who joined this pair game",
  "gamefield": {
    "field": [
      [0, 0, 0],
      [0, 0, 0],
      [0, 0, 0]
    ]
  },
  "state": 0 // gameStatusPair
}
```

- "GET /api/game/player/{PlayerUUID}/info" - get player info  
Response Body:  
```json
{
  "playername": "name..."
}
```

- "GET /api/game/player/{PlayerUUID}/games" - get all players games info  
Response Body:  
```json
[
  {
    "gameuuid": "GameUUID",
    "winneruuid": "...",
    "datecreated": "2026-07-11T09:09:00Z",
    "state": 0 // "StatusesPair"
  }, ...
]
```

- "GET /api/game/leaderboard" - get current leaderboard (Sorted by wins)
QParams:  
"limit" (from 1 to +inf)  
Response Body:  
```json
[
  {
    "uuid": "Player UUID",
    "login": "Player nickname/login",
    "w_ratio": 1.5 // "Player win/(loose+draw) ratio, more = better"
  }, ...
]
```

StatusesSolo:  
- PROGRESS_STATUS = 0
- COMP_WIN_STATUS = 1
- PLYR_WIN_STATUS = 2
- FULL_FLD_STATUS = 3 

StatusesPair:  
- WAIT_OTHER_PLAYER_STATE = 0
-	FirstMOVING_STATE = 1
-	SecondMOVING_STATE = 2
-	DRAW_STATE = 3
-	FirstWIN_STATE = 4
-	SecondWIN_STATE = 5  

Defined cells vals:  
- CELL_EMPTY = 0
- CELL_PLAYER = 1
- CELL_COMPUTER = 2  

# **AuthAPI:**  
- "POST /api/auth/register" - register a new user  
Request Body:  
```json
{
  "login": "Your Login", 
  "password": "Your Password"
} 
``` 
Response Body:  
```json
{
  "status": "success"
}
```

- "POST /api/auth/login" - get your access and refresh jwt  
Request Body:  
```json
{
  "login": "Your Login", 
  "password": "Your Password"
}  
```
Response Body:  
```json
{
  "accessToken": "Your Access token",
  "refreshToken": "Your Refresh token"
}
```

- "POST /api/jwt/update/access" - update only your access jwt  
Request Body:  
```json
{
  "refreshToken": "Your Refresh token"
}  
```
Response Body:  
```json
{
  "accessToken": "Your New Access token",
  "refreshToken": "Your Refresh token"
}
```

- "POST /api/jwt/update/refresh" - update all your jwt  
Request Body:  
```json
{
  "refreshToken": "Your Refresh token"
}  
```
Response Body:  
```json
{
  "accessToken": "Your New Access token",
  "refreshToken": "Your New Refresh token"
}
```