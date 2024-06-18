# versionControlSystem

1. help

2. commit
  descr might be "descr"
  2.1. commit -d <descr> -a <author>        create commit
  
3. branch
  3.1. branch                              список веток
  3.2. branch master                       информация о ветке (10 коммитов)
    3.2.1. -a                              все коммиты
    3.2.2. -v                              подробный вывод
    3.3.3. -c 10                           показать n коммитов

4. checkout
  4.1. checkout <branch>                   переключить ветку
    4.1.1. -b                              создать новую ветку

5. diffs
  5.1. diffs
  5.2. diffs <commithash>
  5.3. diffs <commit1Hash> <commit2Hash>

6. show
  6.1 show <hash>                           показать объект
