package main

import "fmt"

func main() {

	content := "Haskell (/hæskəl/[25]) is a general-purpose, statically typed, purely functional programming language with type inference and lazy evaluation.[26][27] Haskell pioneered several programming language features such as type classes, which enable type-safe operator overloading, and monadic input/output (IO). It is named after logician Haskell Curry.[1] Haskell's main implementation is the Glasgow Haskell Compiler (GHC). \n\n Haskell's semantics are historically based on those of the Miranda programming language, which served to focus the efforts of the initial Haskell working group.[28] The last formal specification of the language was made in July 2010, while the development of GHC continues to expand Haskell via language extensions. \n\n Haskell is used in academia and industry.[29][30][31] As of May 2021, Haskell was the 28th most popular programming language by Google searches for tutorials,[32] and made up less than 1 percent of active users on the GitHub source code repository.[33] Haskell features lazy evaluation, lambda expressions, pattern matching, list comprehension, type classes and type polymorphism. It is a purely functional programming language, which means that functions generally have no side effects. A distinct construct exists to represent side effects, orthogonal to the type of functions. A pure function can return a side effect that is subsequently executed, modeling the impure functions of other languages.\n\n Haskell has a strong, static type system based on Hindley–Milner type inference. Its principal innovation in this area is type classes, originally conceived as a principled way to add overloading to the language,[41] but since finding many more uses.[42] \n\n The construct that represents side effects is an example of a monad: a general framework which can model various computations such as error handling, nondeterminism, parsing and software transactional memory. They are defined as ordinary datatypes, but Haskell provides some syntactic sugar for their use. \n\n Haskell has an open, published specification,[27] and multiple implementations exist. Its main implementation, the Glasgow Haskell Compiler (GHC), is both an interpreter and native-code compiler that runs on most platforms. GHC is noted for its rich type system incorporating recent innovations such as generalized algebraic data types and type families. The Computer Language Benchmarks Game also highlights its high-performance implementation of concurrency and parallelism.[43] \n\n An active, growing community exists around the language, and more than 5,400 third-party open-source libraries and tools are available in the online package repository Hackage.[44]"

	stringKeywords := getStringKeywords(content, 5)

	fmt.Println("String keywords: ", stringKeywords)

	fileKeywords := getFileKeywords("./data/sample.txt", 5) // use term frequency from file

	fmt.Println("File keywords: ", fileKeywords)

	// get keywords for files based on tf-idf
	//filepaths := []string{"./data/sample.txt", "./data/sample2.txt"}

}
