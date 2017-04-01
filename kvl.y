%{
package logtee
%}

%union {
    entries []fieldEntry
    v string
}

%token tFieldName
%token tEquals
%token tEscapeValue
%token tUnescapeValue
%%

line
    : /* empty */
    | tFieldName tEquals tEscapeValue line {$$.entries = append($4.entries, fieldEntry{$1.v, $3.v}); yylex.(*Lexer).parseResult = $$;}
    | tFieldName tEquals tUnescapeValue line {$$.entries = append($4.entries, fieldEntry{$1.v, $3.v}); yylex.(*Lexer).parseResult = $$;}
    | tFieldName tEquals tFieldName line {$$.entries = append($4.entries, fieldEntry{$1.v, $3.v}); yylex.(*Lexer).parseResult = $$;}
    | tFieldName tEquals line {$$.entries = append($3.entries, fieldEntry{$1.v, ""}); yylex.(*Lexer).parseResult = $$;}
    ;

%%