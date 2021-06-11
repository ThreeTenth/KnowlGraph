// app_home.html 
 const app_home = `<div class="app">
  <nav class="nav">
    <div class="nav-start">
      <router-link class="nav-a" to="/">🏠 {{ i18n.Home }}</router-link>
    </div>
    <dropdown>
      <button class="nav-img">
        <img alt="" v-bind:src="profilePicture">
      </button>
      <div class="drop-content">
        <div class="flex">
          <router-link class="flex-1" v-bind:to="{ name: 'newArticle', query: { status: 'public' }}">{{
            i18n.NewArticle }}</router-link>
          <router-link v-bind:to="{ path: '/new/article', query: { status: 'private' }}">🔒</router-link>
        </div>
        <router-link to="/drafts">{{ i18n.Drafts }}</router-link>
        <router-link to="/archive">{{ i18n.Archive }}</router-link>
        <hr />
        <router-link to="/settings/personalize">{{ i18n.Personalize }}</router-link>
        <router-link to="/about">{{ i18n.About }}</router-link>
        <hr />
        <a href="/signout">{{ i18n.Logout }}</a>
      </div>
    </dropdown>
  </nav>
  <router-view></router-view>
  <vs-dialog overflow-hidden v-model="vote.has">
    <div style="margin: 16px 10px;" v-if="vote.has">
      <h1>{{ i18n.DiscoverANewArticle }}</h1>
      <h2 v-if="vote.ras.edges.version.title">{{ vote.ras.edges.version.title }}</h2>
      <p>{{ vote.ras.edges.version.gist }}</p>
      <hr />
      <div v-html="md2html(vote.ras.edges.version.edges.content.body)"></div>
      <div class="article-keywords" v-if="vote.ras.edges.version.edges.keywords">
        <span v-for="(keyword, i) in vote.ras.edges.version.edges.keywords">{{ keyword.name }}</span>
      </div>
      <hr />
      <div class="m-t22">是否同意该文章发布？</div>
      <div class="m-t10">
        <button class="nav-button" v-on:click="onAllow">同意发布</button>
        <button class="nav-button" v-on:click="onRejecte">反对发布</button>
      </div>
    </div>
  </vs-dialog>
</div>`

// app_index.html 
 const app_index = `<div class="app">
  <div class="nav">
    <div class="nav-start">
      <router-link class="nav-a" to="/">🏠 {{ i18n.Home }}</router-link>
    </div>
    <div>
      <button class="nav-button" v-on:click="onGitHubOAuth">{{ i18n.LoginWithGitHubOAuth }}</button>
    </div>
  </div>
  <router-view class="view"></router-view>
</div>`

// com_autocomplet.html 
 const com_autocomplet = `<div class="auto-completion" v-bind:class="pos">
  <input autofocus type="text" v-model="input" class="p-2-4" v-on:keydown="onHotKey" v-on:keyup="onChanged"
    v-bind:placeholder="placeholder" ref="autocomplete">
  <ul class="selection" v-bind:class="items.length ? 'block' : ''" ref="selection">
    <li v-for="(item, i) in items" v-bind:class="selected == i ? 'selected' : ''" v-on:click.stop="onSelect(item)"
      v-bind:key="item.id">{{ item.name }}</li>
  </ul>
</div>`

// com_dropdown.html 
 const com_dropdown = `<div class="drop" v-bind:class="[{block: seen}, pos]" v-on:click="toggle" ref="dropdown">
  <slot></slot>
</div>`

// com_modal.html 
 const com_modal = `<span>
  <transition name="vs-dialog">
    <div v-on:click.stop="close" class="modal" v-if="seen" ref="dialog-content">
      <div class="dialog" v-on:click.stop>
        <button class="close" v-on:click="close">✖</button>
        <div class="dialog__content notFooter">
        </div>
      </div>
    </div>
  </transition>
</span>`

// com_new_node.html 
 const com_new_node = `<div class="new-node">
  <h3>创建新节点</h3>
  <ul class="words" style="margin-bottom: 15px;">
    <li>></li>
    <template v-if="selected">
      <li>
        {{ selected.edges.node.edges.word.name }}
      </li>
      <li>></li>
    </template>
  </ul>
  <div class="nodes">
    <button class="option" v-bind:class="selected == archive ? 'selected' : ''" v-for="(archive, i) in user.archives"
      v-on:click="onSelectPrev(archive)">📁 {{ archive.edges.node.edges.word.name }}</button>
  </div>
  <div class="input">
    <label>新建节点名：</label>
    <div>
      <autocomplet v-on:item="getWordName" v-on:select="onSelectWord" v-bind:source="user.words"
        v-bind:placeholder="i18n.AddAWord"></autocomplet>
    </div>
  </div>
  <hr>
  <div class="options">
    <button class="menu-button">确定</button>
  </div>
</div>`

// com_the_time.html 
 const com_the_time = `<time v-bind:title="fulltime">{{ timetonow }}</time>`

// com_words_input.html 
 const com_words_input = `<div clss="word-layout">
  <div class="word-input">
    <ul class="words">
      <li class="word" v-for="(word, i) in values"
        v-bind:class="i == deleteMark ? 'del' : ''"
        v-on:click="onRemoveWord(word)">{{ word }} X
      </li>
      <li class="new-word">
        <input v-model="input" type="text" v-on:keyup="onKeyup" v-on:keyup.enter="onAdd" v-on:keyup.188="onAdd"
          v-on:keyup.delete="onDel" v-on:keyup.up="onMoveUp" v-on:keyup.down="onMoveDown"
          v-bind:placeholder="i18n.AddAWord"></input>
      </li>
    </ul>
  </div>
  <div class="word-autocomplete" v-if="auto">
    <ul class="">
      <li class="" v-for="item in items">
        <div>{{ item }}</div>
      </li>
    </ul>
  </div>
</div>`

// fgm_about.html 
 const fgm_about = `<h2>{{ i18n.About }}</h2>`

// fgm_archive.html 
 const fgm_archive = `<div class="">
  <h2>{{ i18n.Archive }}</h2>
  <div class="article-nav">
    <router-link to="/archive" v-bind:class="status == 'achiveSelf' ? 'selected' : ''">{{ i18n.MyArticles }} 🔒</router-link>
    <router-link to="/archive/star" v-bind:class="status == 'achiveStar' ? 'selected' : ''">{{ i18n.Star }} ⭐</router-link>
    <router-link to="/archive/watch" v-bind:class="status == 'achiveWatch' ? 'selected' : ''">{{ i18n.Watch }} 👀</router-link>
    <router-link to="/archive/readlist" v-bind:class="status == 'achiveReadlist' ? 'selected' : ''">{{ i18n.ReadList }} 🕒</router-link>
  </div>
  <router-view></router-view>
</div>`

// fgm_archive_articles.html 
 const fgm_archive_articles = `<div v-if="!loading" class="m-t22">
  <div v-if="articles.length == 0" style="text-align: center; margin-top: 120px;">
    <div>
      <span style="font-size: 2em; margin: 30px 0;">{{ statusIcon }}</span>
      <h3>There are not any {{status}} articles.</h3>
    </div>
  </div>
  <div v-else class="articles-item" v-for="(article, i) in articles">
    <router-link class="articles-item-link" :to="{ name: 'article', params: { id: article.id, code: article.code }}">
      <div class="articles-item-body">
        <div v-if="article.title">{{ article.title }}</div>
        <div>{{ article.gist }}</div>
      </div>
    </router-link>
    <the-time v-bind:datetime="article.created_at"></the-time>
  </div>
</div>
<div v-else style="text-align: center; margin-top: 120px;">
  <div class="classic-1"></div>
</div>`

// fgm_article.html 
 const fgm_article = `<div class="container zero">
  <div v-if="article">
    <div class="article-nodes">
      <div>
        <dropdown class="left" v-bind:ignore="'true'" v-if="article.nodes.length">
          <button class="menu-button">{{ article.nodes[0].edges.word.name }}</button>
          <div class="drop-content">
            <div v-for="(node, i) in article.nodes">{{ node.edges.word.name }}</div>
          </div>
        </dropdown>
        <dropdown class="left" v-bind:ignore="'true'" v-else>
          <button class="menu-button">💡 {{ i18n.Archive }}</button>
          <div class="drop-content" v-if="user.logined">
            <div v-if="user.archives.length">
              <button v-for="(archive, i) in user.archives">📁 {{ archive.edges.node.edges.word.name }}</button>
              <hr>
            </div>
            <p v-else>无节点 😃</p>
            <button v-on:click="onNewNode">新建一个节点并关注</button>
            <button>关注更多节点</button>
          </div>
          <div class="drop-content" v-else>
            <router-link to="/login">{{ i18n.PleaseLoginFirst }}</router-link>
          </div>
        </dropdown>
      </div>
      <template v-if="article.private">
        <button class="menu-button">🔒</button>
      </template>
      <template v-else>
        <button class="menu-button" v-on:click="onStar">⭐ {{ article.star ? i18n.Unstar : i18n.Star }}</button>
        <button class="menu-button" v-on:click="onWatch">👀 {{ article.watch ? i18n.Unwatch : i18n.Watch }}</button>
        <button class="menu-button" v-on:click="switchLang">🌍 {{ article.lang.name }}</button>
      </template>
    </div>
    <div class="article-nav">
      <a class="selected" href="#">{{ i18n.Body }}</a>
      <a href="#">{{ i18n.Issues }}</a>
      <a href="#">{{ i18n.Notes }}</a>
    </div>
    <article class="article-content" v-html="article.body"></article>
    <div class="article-keywords" v-if="article.keywords">
      <a href="#" v-for="(keyword, i) in article.keywords">#{{keyword.name}}</a>
    </div>
    <div class="bar">
      <div>
        <the-time v-bind:datetime="article.created_at" />
      </div>
      <dropdown>
        <button class="emoji-button">😃</button>
        <div class="drop-content">
          <div class="p-sm">{{ i18n.PickYourReaction }}</div>
          <div class="flex">
            <button class="reaction-button" v-for="(reac, i) in reactions.slice(0, 5)"
              v-on:click="onPickReaction(reac.value)" v-bind:title="i18n[reac.name]">{{ reac.emoji}}</button>
          </div>
          <div class="flex">
            <button class="reaction-button" v-for="(reac, i) in reactions.slice(5, 8)"
              v-on:click="onPickReaction(reac.value)" v-bind:title="i18n[reac.name]">{{ reac.emoji}}</button>
          </div>
        </div>
      </dropdown>
      <dropdown>
        <button class="emoji-button">
          <svg version="1.1" width="24" height="24" viewBox="0 0 24 24">
            <path
              d="M16,12C16,10.9 16.9,10 18,10C19.1,10 20,10.9 20,12C20,13.1 19.1,14 18,14C16.9,14 16,13.1 16,12M10,12C10,10.9 10.9,10 12,10C13.1,10 14,10.9 14,12C14,13.1 13.1,14 12,14C10.9,14 10,13.1 10,12M4,12C4,10.9 4.9,10 6,10C7.1,10 8,10.9 8,12C8,13.1 7.1,14 6,14C4.9,14 4,13.1 4,12Z" />
          </svg>
        </button>
        <div class="drop-content">
          <button>{{ i18n.Histories }}</button>
          <button>{{ i18n.Translate }}</button>
          <button v-on:click="onEditArticle">{{ i18n.Edit }}</button>
        </div>
      </dropdown>
    </div>
    <div class="article-reactions" v-if="article.reactions.length">
      <button class="reaction-button" v-for="(reac, i) in article.reactions" v-on:click="onPickReaction(reac.status)">{{
        emoji[reac.status]}} {{ reac.count }}</button>
    </div>
  </div>
  <div v-else>Loading</div>
  <modal v-model="isNewNode">
    <new-node></new-node>
  </modal>
</div>`

// fgm_draft_histories.html 
 const fgm_draft_histories = `<div class="container">
  <ol>
    <li class="draft-histories-item" v-for="(snapshot, i) in snapshots" v-on:click="onHistory(snapshot)">
      <the-time v-bind:datetime="snapshot.created_at"></the-time>
    </li>
  </ol>
</div>`

// fgm_draft_history.html 
 const fgm_draft_history = `<div class="container">
  <div v-if="snapshot">
    <p class="wrap">{{ snapshot.body }}</p>
    <div class="bar">
      <div></div>
      <dropdown>
        <button class="emoji-button">
          <svg version="1.1" width="24" height="24" viewBox="0 0 24 24">
            <path
              d="M16,12C16,10.9 16.9,10 18,10C19.1,10 20,10.9 20,12C20,13.1 19.1,14 18,14C16.9,14 16,13.1 16,12M10,12C10,10.9 10.9,10 12,10C13.1,10 14,10.9 14,12C14,13.1 13.1,14 12,14C10.9,14 10,13.1 10,12M4,12C4,10.9 4.9,10 6,10C7.1,10 8,10.9 8,12C8,13.1 7.1,14 6,14C4.9,14 4,13.1 4,12Z" />
          </svg>
        </button>
        <div class="drop-content">
          <button v-on:click="onBack">{{ i18n.Back }}</button>
          <button>{{ i18n.Preview }}</button>
          <button>{{ i18n.Compare }}</button>
          <button>{{ i18n.RestoreThisVersion }}</button>
        </div>
      </dropdown>
    </div>
    <the-time v-bind:datetime="snapshot.created_at"></the-time>
  </div>
  <div v-else>
    Loading
  </div>
</div>`

// fgm_home.html 
 const fgm_home = `<div class="container">
  <div v-if="articles">
    <div class="articles-item" v-for="(article, i) in articles">
      <router-link class="articles-item-link" :to="{ name: 'article', params: { id: article.id, code: article.code }}">
        <div class="articles-item-body">
          <div v-if="article.title">{{ article.title }}</div>
          <div>{{ article.gist }}</div>
        </div>
      </router-link>
      <the-time v-bind:datetime="article.created_at"></the-time>
    </div>
  </div>
  <div v-else>Loading</div>
</div>`

// fgm_index.html 
 const fgm_index = `<h3>{{ i18n.Index }}</h3>`

// fgm_login.html 
 const fgm_login = `<div class="modal white">
  <div>
    <h1>{{ i18n.Login }}
      <router-link to="/">KnowlGraph</router-link>
    </h1>
    <button v-on:click="onGitHubOAuth">{{ i18n.LoginWithGitHubOAuth }}</button>
  </div>
</div>`

// fgm_my.html 
 const fgm_my = `<h3>{{ i18n.My }}</h3>`

// fgm_new_article.html 
 const fgm_new_article = `<div class="container">
  <textarea class="draft-area" v-on:keyup="onChanged" v-on:blur="onBlur" v-model="body"></textarea>
  <div class="bar">
    <div></div>
    <button class="menu-button" v-on:click="onPublish">{{ i18n.Publish }}<span v-if="status=='private'">
        🔒</span></button>
    <dropdown>
      <button class="emoji-button">
        <svg version="1.1" width="24" height="24" viewBox="0 0 24 24">
          <path
            d="M16,12C16,10.9 16.9,10 18,10C19.1,10 20,10.9 20,12C20,13.1 19.1,14 18,14C16.9,14 16,13.1 16,12M10,12C10,10.9 10.9,10 12,10C13.1,10 14,10.9 14,12C14,13.1 13.1,14 12,14C10.9,14 10,13.1 10,12M4,12C4,10.9 4.9,10 6,10C7.1,10 8,10.9 8,12C8,13.1 7.1,14 6,14C4.9,14 4,13.1 4,12Z" />
        </svg>
      </button>
      <div class="drop-content">
        <button v-on:click="onHistories">{{ i18n.Histories }}</button>
      </div>
    </dropdown>
  </div>
</div>`

// fgm_personalize.html 
 const fgm_personalize = `<div>
  <h2>{{ i18n.Personalize }}</h2>
  <h3>{{ currentLang }}</h3>
  <div class="center grid">
    <vs-row>
      <vs-col v-for="(lang, i) in languages" :key="lang.code" vs-type="flex" vs-justify="center" vs-align="center" w="3"
        v-bind:dir="lang.direction">
        <vs-button v-on:click="onSelectUserLang(lang)">{{ lang.name }}</vs-button>
      </vs-col>
    </vs-row>
  </div>
</div>`

// fgm_publish_article.html 
 const fgm_publish_article = `<div class="container">
  <h1>{{ i18n.Publish }}</h1>
  <div>
    <label for="title">{{ i18n.Cover }}</label>
    <input type="file" name="cover" class="publish-input" v-on:change="cover" v-bind:placeholder="i18n.Cover"></input>
    <label for="title">{{ i18n.Title }}</label>
    <input type="text" name="title" class="publish-input" v-model="title" v-bind:placeholder="i18n.Title"></input>
    <label for="gist">{{ i18n.Gist }}</label>
    <textarea type="text" name="gist" class="publish-input" v-model="gist" v-bind:placeholder="i18n.Gist"></textarea>
    <label for="gist">{{ i18n.VersionName }}</label>
    <input type="text" name="versionName" class="publish-input" v-model="versionName"
      v-bind:placeholder="i18n.VersionName"></input>
    <label for="gist">{{ i18n.Comment }}</label>
    <textarea type="text" name="comment" class="publish-input" v-model="comment"
      v-bind:placeholder="i18n.Comment"></textarea>
    <label for="gist">{{ i18n.ChooseALanguage }}</label>
    <select id="userLangSelect" name="lang" class="publish-input" v-model="lang">
      <option v-for="(lang, i) in languages" v-bind:dir="lang.direction" v-bind:value="lang.code">{{ lang.name }}
      </option>
    </select>
    <words-input v-bind:keywords="keywords" v-bind:words="words" />
    <div class="bar">
      <div></div>
      <button class="menu-button" v-on:click="onPublish">{{ i18n.Publish }}<span
          v-if="status=='private'"> 🔒</span></button>
    </div>
  </div>
</div>`

// fgm_user_drafts.html 
 const fgm_user_drafts = `<div class="">
  <div v-if="drafts">
    <h2>{{ i18n.Drafts }}</h2>
    <div class="drafts-item" v-for="(draft, i) in drafts">
      <div v-on:click="onDraft(i)">{{ draft.body }}</div>
      <the-time v-bind:datetime="draft.created_at"></the-time>
    </div>
  </div>
  <div v-else>Loading</div>
</div>`

