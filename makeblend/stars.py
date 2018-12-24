#!/usr/bin/env python3

import sys
import json
#import bpy

def addStar(data):

    if None in data["Quadrants"]:
        x = data["Boundary"]["Center"]["X"]
        y = data["boundary"]["Center"]["Y"]
        w = data["boundary"]["Width"]

        bpy.ops.mesh.primitive_plane_add(size=w, view_align=False, enter_editmode=False, location=(x, y, 0))

    for subtree in data["Quadrants"]:
        if subtree != None:
            addPlane(subtree)

def getVerts(data):

    x = data["star"]["C"]["X"]
    y = data["star"]["C"]["Y"]

    if x != 0 and y != 0:
        localVerts = []
        star = [x, y]
        localVerts.append(star)   

        for subtree in data["Quadrants"]:
            if subtree != None:
                localVerts.append(getVerts(subtree))
         
        if localVerts != None:
            return localVerts

with open("real.json") as f:
    data = json.load(f)

    #addPlane(data[0]["Quadrants"][0])

    verts = []
    verts.append(getVerts(data[0]))
    print(verts)
    
"""
    # Create a mesh and an object
    mesh = bpy.data.meshes.new("0")
    object = bpy.data.objects.new("0",mesh)
    
    # Set the mesh location
    object.location = bpy.context.scene.cursor_location
    bpy.context.scene.objects.linl(object)

    # Create the mesh from the given data points
    mesh.from_pydata(verts,[],[])
    mesh.update(calc_edges=True)
"""
