# outgo

## INTRO
This is a CLI app for self-development. 
It works through resources and playlists.
You can create your own playlists, and add whatever resource you want (add a new resource of your choosing or use one I've added already). 

## What It Looks Like 

![First SS](https://raw.githubusercontent.com/AtillaColak/outgo/refs/heads/main/s1.png)

![Second SS](https://raw.githubusercontent.com/AtillaColak/outgo/refs/heads/main/s2.png)

## HOW TO RUN IT. 
You have 2 options for this: 
1) If you have Go installed on your computer, you can edit/customize the code however you like and compile the executable yourself (which may be technical).
2) Download and run the appropriate release for your OS from the `releases` list. FOR THIS TO WORK, however, you need to run the executable from within the directory (as it needs to access the json files). 
# Technicals and Dev Process 
## Web Scraping
* I scraped all the data you see in the resources.json from various trustable websites, blogposts, forums, and GitHub repos.

## What's Missing
* As I targeted to deploy it fast, I did not spend time on modularizing the code, so most of the functionality is bundled in the `main.go`. This is not ideal, but maybe something I'll fix later.
* Some stuff links for each resource (even though present in the JSON file) are not displayed even with filtering options. I did not choose this because this would mess up the table view.
* The table view should be prettier and the app must be easier to navigate.
* viewing a playlist should be done with the playlist name instead of the id.
* fix the flakiness when unmarshalling the resources list when all of the fields are allowed. 

## Resources
That google sheet: https://docs.google.com/spreadsheets/d/1wganKHEJps87WhFI2O_xyVw-3vkTshmaf665OKczbwc/edit?gid=0#gid=0
