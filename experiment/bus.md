## How Bus Work

- A hierarchy object that has its parent hierarchy object must enable override 
parent in the "User-Defined Auxiliary Sends" section before setting its 
User-Defined Auxiliary Buses.
- User-Defined Auxiliary Send Volume determine the amount of audio signal sent 
to a specific auxiliary bus. If this volume value is set to -96.0, that means 
there's no audio signal being sent.
    - A use case of this is to bypass FX for a hierarchy object.
- The audio signal of an auxiliary bus / bus will propagate to its parent, all 
the way to root bus (not the master audio bus).
- The bus volume determine multiple things. One is determine the amount of 
processed audio signal. It's accumulated starting from the leaf of a hierarchy 
object all the way to the top
