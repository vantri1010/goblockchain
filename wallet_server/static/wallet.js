document.addEventListener('DOMContentLoaded', () => {
    // Lấy thông tin wallet khi trang tải
    async function loadWallet() {
        try {
            const response = await fetch('/wallet', {method: 'GET'});
            if (!response.ok) throw new Error('Failed to load wallet');
            const data = await response.json();
            document.getElementById('public_key').value = data.public_key;
            document.getElementById('private_key').value = data.private_key;
            document.getElementById('blockchain_address').value = data.blockchain_address;
            console.info(data);
        } catch (error) {
            console.error(error);
        }
    }

    // Cập nhật số dư mỗi giây
    async function reloadAmount() {
        try {
            const blockchainAddress = document.getElementById('blockchain_address').value;
            const response = await fetch(`/wallet/amount?blockchain_address=${encodeURIComponent(blockchainAddress)}`, {
                method: 'GET'
            });
            if (!response.ok) throw new Error('Failed to load amount');
            const data = await response.json();
            document.getElementById('wallet_amount').textContent = data.amount;
            console.info(data.amount);
        } catch (error) {
            console.error(error);
        }
    }

    // Xử lý gửi giao dịch
    function handleSendMoney() {
        const sendAmount = parseFloat(document.getElementById('send_amount').value);
        const walletAmount = parseFloat(document.getElementById('wallet_amount').textContent);

        if (isNaN(sendAmount) || sendAmount <= 0) {
            alert('Please enter a valid positive amount.');
            return;
        }

        if (sendAmount > walletAmount) {
            alert('Insufficient funds! You cannot send more than your wallet balance.');
            return;
        }

        const confirmText = 'Are you sure to send?';
        if (!confirm(confirmText)) {
            alert('Canceled');
            return;
        }

        const transactionData = {
            sender_private_key: document.getElementById('private_key').value,
            sender_blockchain_address: document.getElementById('blockchain_address').value,
            recipient_blockchain_address: document.getElementById('recipient_blockchain_address').value,
            sender_public_key: document.getElementById('public_key').value,
            value: sendAmount,
        };

        fetch('/transaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(transactionData),
        })
            .then(response => {
                if (!response.ok) throw new Error('Transaction failed');
                return response.json();
            })
            .then(data => {
                console.info(data);
                alert(data.message === 'fail' ? 'Send fail' : 'Send success');
            })
            .catch(error => {
                console.error(error);
                alert('Send failed');
            });
    }

    // Gọi loadWallet khi trang tải
    loadWallet();

    // Cập nhật số dư mỗi giây
    setInterval(reloadAmount, 5000);

    // Gắn sự kiện click cho nút Send
    document.getElementById('send_money_button').addEventListener('click', handleSendMoney);
});