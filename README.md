# outgo

## INTRO
This is a CLI app for self-development. 
It works through resources and playlists.
You can create your own playlists, and add whatever resource you want (add a new resource of your choosing or use one I've added already). 

## HOW TO RUN IT. 
You have 2 options for this: 
1) If you have Go installed on your computer, you can edit/customize the code however you like and compile the executable yourself (which may be technical).
2) Download and run the appropriate release for your OS from the `releases` list.
# Technicals and Dev Process 
## Web Scraping
* I scraped all the data you see in the resources.json from various trustable websites, blogposts, forums, and GitHub repos.

## What's Missing
* As I targeted to deploy it fast, I did not spend time on modularizing the code, so most of the functionality is bundled in the `main.go`. This is not ideal, but maybe something I'll fix later.
* fetch-updates functionality, which would fetch a new batch of resources is not working yet. I created a publicly accessible Google sheet and tried to add resources there, and have this update functionality to make an HTTP request to get the CSV version of that sheet and update json. For some reason, it did not work (I'm pretty sure because of the inconsistent resource types I used across the application). But I'll fix that later as well.


## Resources
That google sheet: https://docs.google.com/spreadsheets/d/1wganKHEJps87WhFI2O_xyVw-3vkTshmaf665OKczbwc/edit?gid=0#gid=0
