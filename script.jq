#!/usr/bin/jq -f

def olderthan(fmt; days):
    def inner:
        (now - (fmt | gsub("[.].*Z$"; "Z") | fromdateiso8601)) > (days*60*60*24);
    inner;
    
def instances:
    .Reservations[].Instances[];

def instancename:
    (.Tags[] | select(.Key == "Name") | .Value) // "no name";

instances | select(olderthan(.LaunchTime; 2)) | [.LaunchTime, instancename, .InstanceId]
