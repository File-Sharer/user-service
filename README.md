# user-service
Users Microservice

## API Docs
`/api` - base route

**Headers**:
- **`Authorization`**: Bearer `<ACCESS_TOKEN>`

**Designations**:
- **`[AUTH]`** - ***requires** auth*
- **`[PUB]`** - ***doesn't** require auth*

`/auth`:
- **`[PUB]` POST** -> `/signup` - *register*
- **`[PUB]` POST** -> `/signin` - *sign in*
- **`[PUB]` GET** -> `/refresh` - *refresh jwt pair*
- **`[PUB]` GET** -> `/signout` - *sign out*

`/user`:
- **`[AUTH]` GET** -> `/` - *get own info*
- **`[AUTH]` GET** -> `/:<id>` - *get user info with <*id*>*
