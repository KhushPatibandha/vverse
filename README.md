# videoverse take home assignment

Setup:

(1) NixOS: if you are on nixos just type

```
nix-shell
```

in terminal and you are good to go.

(2) Everything else:

- Install Golang anything about 1.20 would work
- Install sqlite3
- ffmpeg, any recent release should work

Clone the project:

```
git clone https://github.com/KhushPatibandha/vverse.git
cd vverse/
```

Start the server:

```
go run cmd/api/main.go
```

Run Test:

```
chmod +x runTests.sh
./runTests.sh
```

Demo Vid:
https://drive.google.com/file/d/1r2xUqpgS7PcqhVBB6KyuuPgQMpadt29r/view?usp=drive_link

Assumptions:

- Since nothing realated to user data was specified so there is no DB for users and hence i have hard coded a single static token to check for authentication.
- Storing videos in a seperate folder and not sqlite db because db could grow very large and incase if the db file is corrupted we can lose all the data and hence storing video files in a seperate folder eventhough files are very small.
- Video's can be of min 5 secs and max 25 secs
- Max Video size is 25 MB with no minimum.
- Expiry of a temp link is 2 mins
- Everytime a video is uploaded a unique Id is assigned that will be used for all the other operations
- Merging two videos results in a new entry to the Db with it's new unique Id.
- But on the other side when a video is trimmed short, it will replace the old video with the newer shorter video.
- Using a temp link to view a video requires no auth.

Resources used (in no particular order):

https://stackoverflow.com/questions/18444194/cutting-multimedia-files-based-on-start-and-end-time-using-ffmpeg

https://www.arj.no/2018/05/18/trimvideo/

https://www.mux.com/articles/stitch-multiple-videos-together-with-ffmpeg

https://youtu.be/4sR77vaEhy8?si=GLO_sEuj8zFMJIXW

https://superuser.com/questions/650291/how-to-get-video-duration-in-seconds

https://ffmpeg.org/faq.html#How-can-I-join-video-files_003f

https://stackoverflow.com/questions/7333232/how-to-concatenate-two-mp4-files-using-ffmpeg

https://shotstack.io/learn/use-ffmpeg-to-concatenate-video/

### Upload a Video

```http
  POST /api/v1/video
```

### Trim a Video

```http
  PUT /api/v1/trim?id=x&s=y&e=z

  Params id(unique int), s(start), e(end)
```

### Merge a video

```http
POST /api/v1/merge?v1=x&v2=y

Params v1(video 1 Id), v2(video 2 Id)
```

### Get temp link

```http
GET /api/v1/link?id=x

Params Id(unique in)
```

### Use temp link to see video

```http
GET /api/v1/uploads/{link}
```
