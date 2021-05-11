# IG-Sniper
An Instagram username taker/claimer - Check desired usernames availability. When it becomes available, original username will be changed to desired username.


## About
This program monitors usernames and checks their availability. It checks each username from `targets.txt` with two methods. The first method is by checking with the web create api, and
the other method is by checking the URL of the username, which can be slightly inaccurate. If the response doesn't show that the username is taken, a post request will be 
sent to the account that has been logged in to update their details (changing username). 

## Preview
![img1](https://i.ibb.co/Lx53N0h/Screenshot-547.png)

## Usage
```
go get "github.com/dlclark/regexp2"
go get "github.com/gookit/color"
```

- Edit `account.json` and put your Instagram login details.
- Edit `targets.txt` and put your desired usernames to check and claim when available
