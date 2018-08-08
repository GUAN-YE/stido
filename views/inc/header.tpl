<!DOCTYPE html>

<html>
    <head>
      <title>独孤影 - {{{.title}}}</title>
      <meta name="viewport" content="width=device-width, initial-scale=1"/>
      <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
      <meta content="独孤影,博客,{{{.keywords}}}" name="keywords" />
      <meta content="{{{.description}}}" name="description" itemprop="description" />
      <meta itemprop="name" content="{{{.articleTitle}}}" />
      <meta itemprop="image" content="{{{.articleImage}}}" />
      <link rel="EditURI" type="application/rsd+xml" title="RSD" href="{{{.host}}}/xmlrpc" />
      <link rel="shortcut icon" href="/favicon.ico" />
      {{{if .inDev}}}
          {{{template "inc/css_dev.tpl" .}}}
      {{{else}}}
          {{{template "inc/css_prod.tpl" .}}}
      {{{end}}}
      <meta name="google-site-verification" content="ohMjRPHv0sKAahvl1H0GC7Dx0-z-zXbMNnWBfxp2PYY" />
      <meta name="baidu-site-verification" content="h3Y69jNgBz" />
    </head>
    <body >
      <div class="main">

          <div class="header">
        		<div class="icons">
        			<a href="http://my.oschina.net/duguying" target="_black"><span title="follow me on oschina" class="icon-osc imgicon"></span></a>
              <a href="https://github.com/duguying" target="_black"><span title="follow me on Github" class="icon-github imgicon"></span></a>
        			<a href="http://weibo.com/duguying2008" target="_black"><span title="find me on Weibo" class="icon-weibo imgicon"></span></a>
        			<a href="http://gplus.to/duguying" target="_black"><span title="find me on g+" class="icon-gplus imgicon"></span></a>
        			<a href="https://twitter.com/duguying" target="_black"><span title="find me on Twitter" class="icon-twitter imgicon"></span></a>
        		</div>
        		<ul class="menu">
        			<li id="about">
                <a>关于</a>
                <div class="drop-menu">
                  <span class="droplist-array-down">◆</span>
                  <ul class="drop-list">
                    <li><a target="_blank" href="">关于博客</a></li>
                    <li><a target="_blank" href="/about/resume">个人简历</a></li>
                    <li><a target="_blank" href="/about/statistics">代码统计</a></li>
                  </ul>
                </div>
              </li>
        			<li>
                <a href="/project">项目</a>
              </li>
        			<li><a href="/list">列表</a></li>
        			<li><a href="/">博文</a></li>
        		</ul>
        		<div class="banner">
        			<a href="/" title="独孤影"><span class="title"><img src="/static/theme/default/img/dgy.svg" alt="独孤影"></span></a>
        		</div>
        		<div class="gap">
              {{{if eq .userIs "admin"}}}
        			<a href="/admin" title="管理页面">
                <img class="gravatar" src="/logo" />
              </a>
              {{{else}}}
              <img class="gravatar" src="/logo" />
              {{{end}}}
        		</div>
        	</div>
