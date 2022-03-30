# Matching

## Change 1: Remove `go`

If the go command is removed everything happens in the same goroutine which means that everything happens in order. That means that the first person will always send a message to the second person and so on.

## Change 2: Declaration of `wg`

If the declaration is changed as specified the separate goroutines will have separate WaitGroups instead of a pointer to a common WaitGroup. This means that there will be a deadlock when the `Seek` routines try to decrement the counter but the counter in the main routine is never reduced.

## Change 3: Remove buffer

If the buffer is removed messages will not be sent unless there is a corresponding receive. This means that the last person to send their message will not be able to do that and the program gets stuck in a deadlock in their case statement.

## Change 4: Remove default case

If the default case is removed the program will be stuck in a deadlock if there is no left over message, in other words if there are an even amount of people