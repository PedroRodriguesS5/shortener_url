const API_URL = 'https://quick-linker.onrender.com';

// Note: Shortened URL redirects are handled directly by the backend API
// Users should access shortened URLs directly at the backend domain

document.addEventListener('DOMContentLoaded', () => {
    const actionBtn = document.getElementById('action-btn');
    const urlInput = document.getElementById('url-input');
    const resultContainer = document.getElementById('result-container');
    const shortenedUrl = document.getElementById('shortened-url');
    const qrcodeImg = document.getElementById('qrcode-img');
    qrcodeImg.crossOrigin = 'anonymous';
    const downloadQrBtn = document.getElementById('download-qr');

    urlInput.focus();

    urlInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            actionBtn.click();
        }
    });

    actionBtn.addEventListener('click', async () => {
        const url = urlInput.value.trim();
        if (!url) {
            showNotification('Please enter a valid URL');
            urlInput.focus();
            return;
        }

        if (!isValidURL(url)) {
            showNotification('Please enter a valid URL (e.g., https://example.com)');
            urlInput.focus();
            return;
        }

        async function fetchAndSetQrCode(url) {
            try {
                const response = await fetch(`${API_URL}/qrcode?url=${encodeURIComponent(url)}`);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const blob = await response.blob();
                const objectURL = URL.createObjectURL(blob);
                qrcodeImg.src = objectURL;
            } catch (error) {
                console.error('Error fetching QR code:', error);
                showNotification('Failed to load QR code.');
            }
        }
        try {
            actionBtn.disabled = true;
            const originalBtnText = actionBtn.textContent;
            actionBtn.textContent = 'Processing...';
            actionBtn.classList.add('loading');
            
            const response = await fetch(`${API_URL}/shorten`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url }),
            });

            actionBtn.disabled = false;
            actionBtn.textContent = originalBtnText;
            actionBtn.classList.remove('loading');

            if (!response.ok) {
                throw new Error('Failed to shorten URL');
            }

            const data = await response.json();
            const shortUrl = data.short_url;

            if (!shortUrl) {
                console.error('Invalid response structure:', data);
                throw new Error('Invalid API response format');
            }

            shortenedUrl.href = shortUrl;
            shortenedUrl.textContent = shortUrl;
            await fetchAndSetQrCode(url);

            resultContainer.classList.remove('hidden');
            setTimeout(() => {
                resultContainer.style.opacity = '1';
                resultContainer.style.transform = 'translateY(0)';
            }, 10);

            addCopyButton(shortUrl);

        } catch (error) {
            console.error('Error:', error);
            showNotification('An error occurred. Please try again.');
            actionBtn.disabled = false;
            actionBtn.textContent = 'Generate';
            actionBtn.classList.remove('loading');
        }
    });


    function isValidURL(string) {
        try {
            new URL(string);
            return true;
        } catch (_) {
            return false;
        }
    }

    function addCopyButton(textToCopy) {
        let copyBtn = document.querySelector('.copy-btn');

        if (!copyBtn) {
            copyBtn = document.createElement('button');
            copyBtn.textContent = 'Copy URL';
            copyBtn.className = 'copy-btn';
            shortenedUrl.insertAdjacentElement('afterend', copyBtn);
        }

        copyBtn.addEventListener('click', async () => {
            try {
                await navigator.clipboard.writeText(textToCopy);
                copyBtn.textContent = 'Copied!';
                showNotification('URL copied to clipboard!');
                setTimeout(() => {
                    copyBtn.textContent = 'Copy URL';
                }, 2000);
            } catch (err) {
                console.error('Failed to copy: ', err);
                showNotification('Failed to copy URL');
            }
        });
    }

    downloadQrBtn.addEventListener('click', () => {
        if (!qrcodeImg.src || qrcodeImg.src === window.location.href) {
            showNotification('No QR code available to download');
            return;
        }

        const originalBtnText = downloadQrBtn.textContent;
        downloadQrBtn.textContent = 'Preparing...';
        downloadQrBtn.disabled = true;

        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');

        const processQRCode = () => {
            try {
                canvas.width = qrcodeImg.naturalWidth;
                canvas.height = qrcodeImg.naturalHeight;

                ctx.fillStyle = '#ffffff';
                ctx.fillRect(0, 0, canvas.width, canvas.height);
                ctx.drawImage(qrcodeImg, 0, 0);

                if ('showSaveFilePicker' in window) {
                    canvas.toBlob(async (blob) => {
                        try {
                            const opts = {
                                suggestedName: 'quicklinker-qr.png',
                                types: [{
                                    description: 'PNG Image',
                                    accept: {'image/png': ['.png']}
                                }]
                            };
                            
                            const handle = await window.showSaveFilePicker(opts);
                            const writable = await handle.createWritable();
                            await writable.write(blob);
                            await writable.close();
                            
                            showNotification('QR code saved successfully!');
                        } catch (err) {
                            if (err.name !== 'AbortError') {
                                console.error('Save error:', err);
                                showNotification('Failed to save QR code');
                            }
                        } finally {
                            downloadQrBtn.textContent = originalBtnText;
                            downloadQrBtn.disabled = false;
                        }
                    }, 'image/png');
                } else {
                    canvas.toBlob(blob => {
                        const url = URL.createObjectURL(blob);
                        const link = document.createElement('a');
                        link.download = 'quicklinker-qr.png';
                        link.href = url;
                        link.click();
                        
                        setTimeout(() => URL.revokeObjectURL(url), 100);
                        showNotification('QR code downloaded successfully!');
                        
                        downloadQrBtn.textContent = originalBtnText;
                        downloadQrBtn.disabled = false;
                    }, 'image/png');
                }
            } catch (err) {
                console.error('QR processing error:', err);
                showNotification('Failed to process QR code');
                downloadQrBtn.textContent = originalBtnText;
                downloadQrBtn.disabled = false;
            }
        };

        if (qrcodeImg.complete && qrcodeImg.naturalHeight !== 0) {
            processQRCode();
        } else {
            qrcodeImg.onload = processQRCode;
        }
    });

    function showNotification(message) {
        const existingNotifications = document.querySelectorAll('.notification');
        existingNotifications.forEach(n => n.remove());

        const notification = document.createElement('div');
        notification.className = 'notification';
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        setTimeout(() => notification.classList.add('show'), 10);
        
        setTimeout(() => {
            notification.classList.remove('show');
            setTimeout(() => notification.remove(), 300);
        }, 3000);
    }
});