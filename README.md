# Description
A basic 3d rendering of a Rubiks cube.
Rotate the cube using your arrow keys
Turn faces of the cube using the keys (F[ront] -> rotate green clockwise, R[ight] -> rotates red clockwise, L[eft] -> orange, B[ack] -> blue, U[p] -> white, D[own] -> yellow)

# How it works:
The struct Side is defined by the 4 edges of its sticker in 3d space and the vector orthogonal to it, as well as a color.
The Rubik's cube struct holds an array of center pieces, each of which contain arrays of their respective corner and edge pieces, which are represented by structs which contain two or three sides (one for each sticker on the piece)
When the a face is rotated, a rotation matrix is applied to the coordinates of all the pieces owned by that side (one piece can be owned by multiple sides). This is broken down into a series of 90, 1 degree rotations, each rendered on each refresh of the screen. At the conclusion of the animation, each face checks whether all the pieces contain sides parallel to that face to update which pieces are owned by each face.
The viewing plane struct holds the coordinate location of the observer as well as a vertical and horizontal vector which form a basis for the plane (which represents the screen). When the arrow keys are pressed, the location of the viewer, and the vertical and horizontal vectors rotate.
The stickers are mapped to the screen by checking whether they are visible using a dot product with the vector from the origin to the viewing plane. They are then ordered by distance from the viewer and rendered closest last. To render them, the project of each edge piece to the viewing plane is found (with respect to the vertical and horizontal vectors) and then the pixels inside the resultant shape are colored in.

# Preview
[![Watch the video](https://img.youtube.com/vi/aBywzsNVLMc/maxresdefault.jpg)](https://youtu.be/aBywzsNVLMc)
