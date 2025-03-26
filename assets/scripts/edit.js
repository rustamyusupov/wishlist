const handleSubmit = async event => {
  event.preventDefault();

  const formData = new FormData(form);
  const url = form.getAttribute('action');
  const body = new URLSearchParams([...formData.entries()]).toString();

  const response = await fetch(url, {
    method: event.submitter.id,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body,
  });

  if (response.ok) {
    window.location.href = '/';
  }
};

const init = () => {
  const form = document.getElementById('edit');

  if (form) {
    form.addEventListener('submit', handleSubmit);
  }
};

document.addEventListener('DOMContentLoaded', init);
