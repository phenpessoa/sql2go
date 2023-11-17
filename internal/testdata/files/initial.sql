-- name: Foo
-- This is a comment that must be ignored
SELECT * FROM foo;

-- name: Bar
SELECT * FROM bar
-- hard coded id for testing
WHERE id = 123;

-- name: Baz
SELECT
    *
FROM
    -- test trimming of whitespaces and comments
    baz
WHERE
    baz = 123 AND
    baz = baz;

-- name: Qux
-- break;
SELECT * FROM qux;

-- name: Quux
SELECT * FROM quux -- this is a comment;
WHERE quux = 123;

-- Random comments
-- will the parser handle these properly?

-- name: Corge
SELECT '--' FROM corge; -- Last comment;

-- name: Grault
SELECT '
-- name: Grault
' FROM grault;

-- name: HardToLex -- this-is-a-comment;
SELECT;

-- name: Empty

-- name: Garply
SELECT 'garply-string-literal' FROM garply;

-- name: Waldo
SELECT "waldo_identifier_1" FROM waldo;

-- name: Fred
SELECT `fred_identifier_2` FROM fred;

-- name: Whatif -- let's see
SELECT * FROM whatif;

-- name: WhatAboutThis
SELECT 'foo--hard--string--literal' FROM whatAboutThis; `foo"_identifier_3'`

-- name: WhatAboutThis2
SELECT

`foo"_identifier_4`

-- name: Plugh
SELECT * FROM plugh

-- name: Xyzzy
SELECT * FROM xyzzy

-- name: Thud
SELECT * FROM thud;
SELECT * FROM thud2;