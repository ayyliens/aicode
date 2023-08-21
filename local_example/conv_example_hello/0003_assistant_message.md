Here is the equivalent "hello world" implementation in Common Lisp using SBCL:

```lisp
(defun hello-world ()
  (format t "Hello, world!"))

(hello-world)
```

In Common Lisp, you define functions using the `defun` form. The `hello-world` function is defined to print "Hello, world!" using the `format` function, with `t` as the first argument indicating that the output should go to the standard output stream. 

To run this code, you can either load it in an interactive Lisp environment like Emacs + SLIME or save the code in a file with a `.lisp` extension and execute it with SBCL. When executed, it will print "Hello, world!" to the console.