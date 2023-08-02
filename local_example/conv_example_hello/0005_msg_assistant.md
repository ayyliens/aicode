Here is the equivalent "hello world" implementation in Haskell:

```haskell
main :: IO ()
main = putStrLn "Hello, world!"
```

In Haskell, the entry point of a program is the `main` function, which has a type signature of `IO ()`. Inside the `main` function, we use the `putStrLn` function to print "Hello, world!" to the console.

To run this code, save it in a file with a `.hs` extension and then compile and execute it using the GHC compiler. Alternatively, you can run it in a Haskell interpreter, such as GHCi, and call the `main` function directly. When executed, it will print "Hello, world!" to the console.