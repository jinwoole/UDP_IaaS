import{a as d,t as b,c as Z}from"../chunks/DldQStsR.js";import{i as tt}from"../chunks/CATt0WfG.js";import{p as et,o as at,a as rt,b as o,f as V,c as ot,g as a,s as $,e as s,r,h as st,i as nt,j as T}from"../chunks/FJRCtc8Y.js";import{e as p,s as v}from"../chunks/CuXJoPns.js";import{i as k}from"../chunks/B2zXuP2w.js";import{e as lt,f as ct}from"../chunks/DpdfpHoe.js";import{c as it,C as dt}from"../chunks/DvOSdxAZ.js";import{s as pt}from"../chunks/CyGVOWoH.js";var ht=b('<button class="bg-gray-600 text-white px-3 py-1 rounded hover:bg-gray-700">Stop</button> <button class="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700">Console</button>',1),vt=b('<button class="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700">Start</button>'),bt=b('<tr class="border-b"><td class="py-2"> </td><td class="py-2"> </td><td class="py-2"> </td><td class="py-2"><span> </span></td><td class="py-2 space-x-2"><!> <button class="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700">Delete</button></td></tr>'),ut=b('<div class="flex items-center justify-center min-h-screen"><h1 class="text-2xl font-semibold text-gray-800">Welcome to UDP IaaS</h1></div> <div class="bg-white rounded-lg shadow p-6"><div class="mb-4"><button class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Create VM</button></div> <table class="w-full"><thead><tr class="border-b"><th class="text-left py-2">Name</th><th class="text-left py-2">Cores</th><th class="text-left py-2">Memory</th><th class="text-left py-2">State</th><th class="text-left py-2">Actions</th></tr></thead><tbody></tbody></table></div> <!>',1);function St(A,O){et(O,!1);let C=T([]),u=T(!1),D;async function l(){try{const e=await fetch("/api/vms");$(C,await e.json())}catch(e){console.error("Failed to load VMs:",e)}}async function q(e){try{await fetch(`/api/vms/${e}/start`,{method:"POST"}),await l()}catch(t){console.error("Failed to start VM:",t)}}async function B(e){try{await fetch(`/api/vms/${e}/stop`,{method:"POST"}),await l()}catch(t){console.error("Failed to stop VM:",t)}}async function L(e){if(confirm(`Are you sure you want to delete ${e}?`))try{await fetch(`/api/vms/${e}`,{method:"DELETE"}),await l()}catch(t){console.error("Failed to delete VM:",t)}}async function U(e){try{const t=await fetch(`/api/vms/${e}/vnc`),{port:c}=await t.json();window.open(`http://${window.location.hostname}:${c}/vnc.html`,"_blank")}catch(t){console.error("Failed to open VNC:",t)}}at(()=>{l(),D=setInterval(l,5e3)}),rt(()=>{clearInterval(D)}),tt();var F=ut(),f=o(V(F),2),m=s(f),W=s(m);r(m);var j=o(m,2),E=o(s(j));lt(E,5,()=>a(C),e=>e.name,(e,t)=>{var c=bt(),y=s(c),H=s(y,!0);r(y);var _=o(y),J=s(_,!0);r(_);var x=o(_),K=s(x);r(x);var w=o(x),g=s(w),Q=s(g,!0);r(g),r(w);var I=o(w),N=s(I);{var R=n=>{var i=ht(),h=V(i),M=o(h,2);p("click",h,()=>B(a(t).name)),p("click",M,()=>U(a(t).name)),d(n,i)},X=n=>{var i=Z(),h=V(i);{var M=S=>{var P=vt();p("click",P,()=>q(a(t).name)),d(S,P)};k(h,S=>{a(t).state==="stopped"&&S(M)},!0)}d(n,i)};k(N,n=>{a(t).state==="running"?n(R):n(X,!1)})}var Y=o(N,2);r(I),r(c),st(n=>{v(H,a(t).name),v(J,a(t).cores),v(K,`${a(t).memory??""} MB`),pt(g,it(n)),v(Q,a(t).state)},[()=>ct(a(t).state)],nt),p("click",Y,()=>L(a(t).name)),d(e,c)}),r(E),r(j),r(f);var z=o(f,2);{var G=e=>{dt(e,{$$events:{close:()=>$(u,!1),created:l}})};k(z,e=>{a(u)&&e(G)})}p("click",W,()=>$(u,!0)),d(A,F),ot()}export{St as component};
