import{a as h,t as S}from"../chunks/DldQStsR.js";import{i as q}from"../chunks/CATt0WfG.js";import{p as M,o as U,c as N,b as n,e as a,s as f,k as P,r as t,h as T,i as A,g as o,j as g}from"../chunks/FJRCtc8Y.js";import{e as w,s as _}from"../chunks/CuXJoPns.js";import{e as B,a as C}from"../chunks/DpdfpHoe.js";var E=S('<tr class="border-b last:border-none hover:bg-gray-50"><td class="py-3 px-4 text-gray-700"> </td><td class="py-3 px-4 text-gray-700"> </td></tr>'),G=S('<div class="min-h-screen bg-gradient-to-r from-gray-200 to-gray-400 flex flex-col items-center justify-center px-4 py-8"><div class="w-full max-w-3xl bg-white rounded-3xl shadow-2xl p-8"><h1 class="text-2xl font-semibold text-gray-800 mb-6 text-center">ISO Management</h1> <form class="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-4 mb-6"><input class="block w-full sm:w-auto text-sm text-gray-700 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:bg-blue-600 file:text-white hover:file:bg-blue-700 cursor-pointer" type="file" accept=".iso" required> <button type="submit" class="bg-blue-600 text-white px-6 py-2 rounded-full hover:bg-blue-800 transition-colors focus:outline-none">Upload</button></form> <div class="overflow-hidden rounded-xl"><table class="min-w-full border-collapse"><thead class="bg-gray-100"><tr><th class="py-3 px-4 text-left text-gray-600">Name</th><th class="py-3 px-4 text-left text-gray-600">Size</th></tr></thead><tbody></tbody></table></div></div></div>');function R(O,j){M(j,!1);let p=g([]),l=g(null);async function m(){try{const e=await fetch("/api/isos");f(p,await e.json())}catch(e){console.error("Failed to load ISOs:",e)}}async function F(e){if(e.preventDefault(),!o(l))return;const r=new FormData;r.append("iso",o(l));try{await fetch("/api/isos",{method:"POST",body:r}),f(l,null),e.target.reset(),await m()}catch(s){console.error("Failed to upload ISO:",s)}}U(m),q();var c=G(),u=a(c),i=n(a(u),2),I=a(i);P(2),t(i);var x=n(i,2),b=a(x),y=n(a(b));B(y,5,()=>o(p),e=>e.name,(e,r)=>{var s=E(),d=a(s),z=a(d,!0);t(d);var v=n(d),D=a(v,!0);t(v),t(s),T(k=>{_(z,o(r).name),_(D,k)},[()=>C(o(r).size)],A),h(e,s)}),t(y),t(b),t(x),t(u),t(c),w("change",I,e=>f(l,e.target.files[0])),w("submit",i,F),h(O,c),N()}export{R as component};
