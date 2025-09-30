document.addEventListener('DOMContentLoaded', function(){
  const form = document.getElementById('contactForm');
  if(!form) return;
  const resultEl = document.getElementById('formResult');
  form.addEventListener('submit', async function(e){
    e.preventDefault();
    resultEl.textContent = '正在发送...';
    try{
      const payload = Object.fromEntries(new FormData(form).entries());
      const res = await fetch(form.action, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      const json = await res.json().catch(() => ({}));
      if(res.ok){
        resultEl.style.color = 'green';
        resultEl.textContent = json.message || '感谢您的留言，我们已收到并会尽快联系您。';
        form.reset();
      } else {
        resultEl.style.color = 'red';
        resultEl.textContent = json.error || '发送失败，请稍后重试。';
      }
    } catch(err){
      resultEl.style.color = 'red';
      resultEl.textContent = '无法发送，请检查网络或稍后再试。';
      console.error(err);
    }
  });
});
