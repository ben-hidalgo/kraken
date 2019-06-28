select sleep(10), case when @foo = "bar" then 1 else (select table_name from information_schema.tables) end;
