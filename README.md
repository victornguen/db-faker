# DB faker

> :warning: **This project is in work in progress.**


This is a tool for generating fake data for testing purposes.

## Installation

```bash
go install github.com/victornguen/db-faker
```

## Usage example

Let's assume you have a simple database hosted in your local machine with the following schema:

```sql
-- 1. Users table
CREATE TABLE users
(
    user_id    SERIAL PRIMARY KEY,
    username   VARCHAR(50)         NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Products table
CREATE TABLE products
(
    product_id     SERIAL PRIMARY KEY,
    product_name   VARCHAR(100)   NOT NULL,
    price          DECIMAL(10, 2) NOT NULL,
    stock_quantity INT            NOT NULL CHECK (stock_quantity >= 0)
);

-- 3. Orders table (depends on users)
CREATE TABLE orders
(
    order_id   SERIAL PRIMARY KEY,
    user_id    INT  NOT NULL,
    order_date DATE NOT NULL,
    status     VARCHAR(20) DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

-- 4. Order items table (depends on orders and products)
CREATE TABLE order_items
(
    order_item_id SERIAL PRIMARY KEY,
    order_id      INT            NOT NULL,
    product_id    INT            NOT NULL,
    quantity      INT            NOT NULL CHECK (quantity > 0),
    unit_price    DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (order_id),
    FOREIGN KEY (product_id) REFERENCES products (product_id)
);
```

Optionally, you can define generation rules in a YAML file like this:

`rules.yaml`:
```yaml
rules:
  orders:
    num: 30
    columns:
      amount: integer(1, 100)
      date: timestamp
      status: oneof[new%20, paid%30, shipped%50]
  users:
    num: 100
  products:
    num: 1000
    columns: {}
  order_items:
    num: 1000
```


You can generate fake data for this schema using the following command:

```bash
db-faker generate --user postgres --password postgres --db my_database_name --rules ./rules.yaml
```

## Generating rules

The rules file is a YAML file with the following structure:

```yaml
rules:
  table_name:
    num: 1000
    columns:
      column_name: rule
      ...
```

The `num` field is the number of rows to generate for the table. The `columns` field is a map where the key is the column name and the value is the rule to generate the data for that column.

The available rules are:
- `oneof[option1%50, option2%50]`: Select one of the options with the given probability. The probability is a number after the `%` symbol. The sum of all probabilities must be 100.
- `constant[value]`: Always use the same value.
- `int|integer`, `int|integer(lower, upper)`: Generate a random integer number.
- `sentence|text`, `sentence|text(n)`: Generate a random sentence with `n` chars length(if set).
- `firstname|name`: random first name.
- `lastname`: random last name.
- `email`: random email.
- `username`:  random username.
- `currency`: random currency code.
- `ccnumber`: random credit card number.
- `cctype`: random credit card type.
- `country`: random country name.
- `city`: random city name.
- `address`: random address.
- `state`: random state name.
- `postalcode`: random postal code.
- `latitude|lat`: random latitude.
- `longitude|lon`: random longitude.
- `phone`: random phone number.
- `date`: random date.
- `dayofweek`: random day of the week.
- `month`: random month.
- `year`: random year.
- `time`: random time.
- `datetime|timestamp`: random timestamp.
- `bloodtype`: random blood type.
- `bloodrhfactor`: random blood Rh factor.
- `bloodgroup`: random blood group.
- `paragraph`: random text paragraph.
- `ipv4`: random IPv4 address.
- `ipv6`: random IPv6 address.
- `mac`: random MAC address.
- `url`: random URL.
- `useragent`: random user agent.

