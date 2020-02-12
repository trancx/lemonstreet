# Article Service
Micro-service focusing on dealing the article. Post, Get or
Delete an article. More articles retrieve of users is assosicated with search engine.


## HTTP
1. Get  /:user/:article
>1 anonymous access

Body is empty with no author infos and user infos. So article module need to provide
an interface FROM_NAME_TO_ARTICLE, and validate the info. 1. search title 2. compare user

>2 authentic access

Body contains user infos and article base info, easily retrieved by the server. 


2\. POST /:author/:title

Must be authentic or denial.

user_base_info, article_content uploaded and server generate the aid and article_base_info
into DB, of cause, cache.

##RPC

