import { blob, integer, real, sqliteTable, text } from 'drizzle-orm/sqlite-core';

const createdAt = () =>
	integer('created_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date());

export const users = sqliteTable('users', {
	id: integer('id').primaryKey({ autoIncrement: true }),
	username: text('username').notNull().unique(),
	createdAt: createdAt()
});

export const credentials = sqliteTable('credentials', {
	id: text('id').primaryKey(),
	userId: integer('user_id')
		.notNull()
		.references(() => users.id, { onDelete: 'cascade' }),
	publicKey: blob('public_key', { mode: 'buffer' }).notNull(),
	counter: integer('counter').notNull().default(0),
	transports: text('transports'),
	createdAt: createdAt()
});

export const categories = sqliteTable('categories', {
	id: integer('id').primaryKey({ autoIncrement: true }),
	name: text('name').notNull().unique(),
	createdAt: createdAt()
});

export const currencies = sqliteTable('currencies', {
	id: integer('id').primaryKey({ autoIncrement: true }),
	code: text('code').notNull().unique(),
	symbol: text('symbol').notNull(),
	createdAt: createdAt()
});

export const wishes = sqliteTable('wishes', {
	id: integer('id').primaryKey({ autoIncrement: true }),
	userId: integer('user_id')
		.notNull()
		.references(() => users.id, { onDelete: 'cascade' }),
	categoryId: integer('category_id')
		.notNull()
		.references(() => categories.id),
	name: text('name').notNull(),
	link: text('link').notNull(),
	sort: integer('sort').notNull().default(0),
	createdAt: createdAt()
});

export const prices = sqliteTable('prices', {
	id: integer('id').primaryKey({ autoIncrement: true }),
	wishId: integer('wish_id')
		.notNull()
		.references(() => wishes.id, { onDelete: 'cascade' }),
	amount: real('amount').notNull(),
	currencyId: integer('currency_id')
		.notNull()
		.references(() => currencies.id),
	createdAt: createdAt()
});
