<?php

namespace Database\Seeders;

use App\Models\Product;
use Illuminate\Database\Seeder;

class ProductSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        $products = [
            [
                'name' => 'Laptop Pro 15',
                'description' => 'High-performance laptop with 16GB RAM and 512GB SSD',
                'price' => 1299.99,
                'quantity' => 25,
                'category' => 'Electronics'
            ],
            [
                'name' => 'Wireless Headphones',
                'description' => 'Noise-cancelling wireless headphones with 30h battery life',
                'price' => 199.99,
                'quantity' => 50,
                'category' => 'Electronics'
            ],
            [
                'name' => 'Coffee Maker',
                'description' => 'Programmable coffee maker with 12-cup capacity',
                'price' => 89.99,
                'quantity' => 15,
                'category' => 'Home & Kitchen'
            ],
            [
                'name' => 'Running Shoes',
                'description' => 'Lightweight running shoes with excellent cushioning',
                'price' => 129.99,
                'quantity' => 30,
                'category' => 'Sports & Outdoors'
            ],
            [
                'name' => 'Smartphone',
                'description' => 'Latest smartphone with 128GB storage and triple camera',
                'price' => 699.99,
                'quantity' => 40,
                'category' => 'Electronics'
            ]
        ];

        foreach ($products as $product) {
            Product::create($product);
        }
    }
}
