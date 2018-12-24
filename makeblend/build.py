#!/usr/bin/env python3

import sys
import json
import bpy

def addPlane(data):

    if None in data["Quadrants"]:
        x = data["boundary"]["Center"]["X"]
        y = data["boundary"]["Center"]["Y"]
        w = data["boundary"]["Width"]

        bpy.ops.mesh.primitive_plane_add(size=w, view_align=False, enter_editmode=False, location=(x, y, 0))

    for subtree in data["Quadrants"]:
        if subtree != None:
            addPlane(subtree)

with open("2500.json") as f:
    data = json.load(f)

    addPlane(data[0]["Quadrants"][0])
