import{E as i,a as l,Y as p,o as m,H as _}from"./main.da501797.js";import{i as d}from"./rpc.01fe265c.js";const c="fanhan-ui";var u={set(o,t,e=0){const s=`${c}:${o}`,r={value:t,_expired:e?Date.now()+e*1e3:!1};window.sessionStorage.setItem(s,JSON.stringify(r))},get(o){const t=`${c}:${o}`,e=JSON.parse(window.sessionStorage.getItem(t));return(e==null?void 0:e._expired)&&Date.now()>e._expired?(this.remove(o),null):e?e.value:null},remove(o){const t=`${c}:${o}`;window.sessionStorage.removeItem(t)}};const w={__name:"Callback",setup(o){const t=i(),e=l(),s=p(),r=t.query.token,n=decodeURIComponent(t.query.redirect);return console.log("sso callback",r,n),r?(s.user={token:r},d.post("/account/myinfo").then(a=>{s.user=a,u.set("user",a),e.replace({path:n})})):(s.user=null,u.remove("user"),e.replace({path:n})),(a,f)=>(m(),_("div"))}};export{w as default};
