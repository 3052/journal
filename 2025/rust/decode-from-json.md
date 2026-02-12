# Decode from JSON

## dynamic

~~~
D:\rust\bin\cargo.exe add --no-default-features rust_json
1       rust_json v0.1.4
PS D:\Desktop\_\vendor\rust_json> scc.exe --exclude-dir tests
Language                 Files     Lines   Blanks  Comments     Code
Rust                         7       858       76        77      705

# cargo.exe add --no-default-features humphrey_json
1       humphrey_json v0.2.2
# scc.exe --exclude-dir tests
Language                 Files     Lines   Blanks  Comments     Code
Rust                         8      1293      132       260      901

# cargo.exe add --no-default-features sj
1       sj v0.23.0
# scc.exe --exclude-dir tests
Language                 Files     Lines   Blanks  Comments     Code
Rust                        24      3450      487       647     2316

# cargo.exe add --no-default-features --features std serde_json
1       serde_json v1.0.85
4       serde v1.0.144
~~~

## static

~~~
D:\rust\bin\cargo.exe add --no-default-features rust_json
D:\rust\bin\cargo.exe add --no-default-features rust_json_derive
1       rust_json v0.1.4
2       rust_json_derive v0.1.1 (proc-macro)
6       syn v1.0.100
PS D:\Desktop\_\vendor> scc.exe --exclude-dir tests
Language                 Files     Lines   Blanks  Comments     Code
Rust                        83     54152     2871      7176    44105

# cargo.exe add --no-default-features --features derive humphrey_json
1       humphrey_json v0.2.2
6       syn v1.0.100
# scc.exe --exclude-dir tests
Language                 Files     Lines   Blanks  Comments     Code
Rust                        85     54657     2966      7377    44314

# cargo.exe add --no-default-features --features serde_derive serde
# cargo.exe add --no-default-features --features std serde_json
1       serde v1.0.144
7       serde_json v1.0.85
9       ryu v1.0.11
# scc.exe --exclude-dir tests
Language                 Files     Lines   Blanks  Comments     Code
Rust                       173     97793     7041     15668    75084
~~~
