// Retrieve the basket from localStorage or initialize it as an empty array
let basket = JSON.parse(localStorage.getItem('basket')) || [];

// Function to calculate total price
function calculateTotalPrice() {
    return basket.reduce((total, item) => total + item.price * item.quantity, 0);
}

// Function to update the total price in the UI
function updateTotalPrice() {
    const totalPriceElement = document.getElementById("total-price");
    totalPriceElement.textContent = `$${calculateTotalPrice().toFixed(2)}`;
}

// Function to render the basket
function renderBasket() {
    const basketItems = document.getElementById("basket-items");
    basketItems.innerHTML = ""; // Clear existing items

    console.log("Rendering Basket: ", basket); // Debugging: log current basket data

    // Render each item in the basket
    basket.forEach((item, index) => {
        const basketItem = document.createElement("div");
        basketItem.classList.add("containers-baskets");
        basketItem.innerHTML = `
            <p class="burger-name">${item.name}</p>
            <p class="burger-price">$${item.price.toFixed(2)}</p>
            <div class="quantity-controls">
                <button class="decrease" data-index="${index}">-</button>
                <span class="quantity">${item.quantity}</span>
                <button class="increase" data-index="${index}">+</button>
            </div>
            <button class="remove" data-index="${index}">x</button>
        `;
        basketItems.appendChild(basketItem);
    });

    updateTotalPrice(); // Update the total price whenever the basket changes
}

// Function to add an item to the basket
function addToBasket(item) {
    const existingItem = basket.find(basketItem => basketItem.id === item.id);

    if (existingItem) {
        existingItem.quantity++; // Increment quantity if item exists
    } else {
        basket.push({ ...item, quantity: 1 }); // Add new item
    }

    // Save the updated basket to localStorage
    localStorage.setItem('basket', JSON.stringify(basket));

    console.log("Basket updated and saved to localStorage: ", basket); // Debugging: log updated basket
    renderBasket(); // Update the basket display
}

// Handle increase, decrease, and remove button clicks
document.getElementById("basket-items").addEventListener("click", event => {
    const index = parseInt(event.target.getAttribute("data-index"));

    if (event.target.classList.contains("increase")) {
        basket[index].quantity++;
    } else if (event.target.classList.contains("decrease")) {
        basket[index].quantity--;
        if (basket[index].quantity === 0) {
            basket.splice(index, 1); // Remove item if quantity is 0
        }
    } else if (event.target.classList.contains("remove")) {
        basket.splice(index, 1); // Remove item
    }

    // Save the updated basket to localStorage
    localStorage.setItem('basket', JSON.stringify(basket));

    console.log("Basket updated after click and saved to localStorage: ", basket); // Debugging: log updated basket
    renderBasket(); // Update the basket display
});

// Attach event listeners to all "Add" buttons
document.querySelectorAll(".add-button").forEach(button => {
    button.addEventListener("click", event => {
        const id = event.target.getAttribute("data-id");
        const name = event.target.getAttribute("data-name");
        const price = parseFloat(event.target.getAttribute("data-price"));

        addToBasket({ id, name, price });
    });
});

// Initialize the basket on page load
window.addEventListener("load", () => {
    renderBasket(); // This will run once the page is fully loaded
    console.log("Basket data loaded from localStorage: ", basket); // Debugging: log loaded basket
});

// ********* Navbar *********
document.addEventListener('DOMContentLoaded', function () {
    const currentPath = window.location.pathname;
    console.log("Current Path: ", currentPath); // Debugging line to check the path
    
    // Check the current path and add 'active-link' to the corresponding item
    if (currentPath === "/" || currentPath === "/index") { // Handle homepage
        document.getElementById("burgers-link").classList.add("active-link");
    } else if (currentPath === "/fishes") {
        document.getElementById("fishes-link").classList.add("active-link");
    } else if (currentPath === "/drinks") {
        document.getElementById("drinks-link").classList.add("active-link");
    }
});



function handleBuyNow() {
    // Ask the user to enter their delivery address with a default detailed address (example)
    const deliveryAddress = prompt("Please enter your delivery address:", "123, Alisher Navoi Street, Apartment 7, Mirzo Ulugbek District, Tashkent");

    if (deliveryAddress) {
        const order = {
            order_products: basket.map(item => ({
                id: parseInt(item.id), // Ensure 'id' is an integer
                quantity: item.quantity // 'quantity' should be a number
            })),
            delivery_address: deliveryAddress
        };

        console.log("Order to send:", order); // Debugging: Log order data to the console

        fetch('/order/purchase', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(order) // Send only relevant data
        })
        .then(response => {
            if (response.redirected) {
                // If redirected, navigate to the new URL
                window.location.href = response.url;
            }
        })
        .then(response => {
            // Check if the response is OK (status code 2xx)
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json(); // Parse the JSON from the response
        })
        .then(data => {
            console.log("Server response:", data); // Log the server response to check

            if (data.success) {
                alert("Order placed successfully!");
                localStorage.removeItem('basket'); // Clear the basket
                basket = []; // Empty the basket array
                renderBasket(); // Update the basket display
            } else {
                console.error("Error placing order:", data.error || 'Unknown error');
                // Here we are not showing the alert anymore, we simply log the error
            }
        })
       
    } else {
        alert("Delivery address is required.");
    }
}



// Attach event listener to the "Buy Now" button
document.getElementById("buy-button").addEventListener("click", handleBuyNow);

// Example to simulate adding items to the basket
function addItemToBasket(id, name, price) {
    const existingItem = basket.find(item => item.id === id);
    if (existingItem) {
        existingItem.quantity++;
    } else {
        basket.push({ id, name, price, quantity: 1 });
    }

    // Save basket to localStorage
    localStorage.setItem('basket', JSON.stringify(basket));
    renderBasket();
}

// Example to simulate rendering the basket on page load
document.addEventListener('DOMContentLoaded', () => {
    renderBasket();
});
