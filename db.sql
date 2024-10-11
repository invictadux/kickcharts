CREATE TABLE IF NOT EXISTS channels(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    banner TEXT NOT NULL,
    picture TEXT NOT NULL,
    is_banned BOOLEAN NOT NULL,
    language TEXT NOT NULL,
    live BOOLEAN NOT NULL,
    live_viewers INTEGER NOT NULL,
    followers_count INTEGER NOT NULL,
    peak_viewers INTEGER NOT NULL,
    description TEXT NOT NULL,
    discord TEXT NOT NULL,
    facebook TEXT NOT NULL,
    instagram TEXT NOT NULL,
    tiktok TEXT NOT NULL,
    twitter TEXT NOT NULL,
    youtube TEXT NOT NULL
);

CREATE INDEX channels_idx_1 ON channels(slug);
CREATE INDEX channels_idx_2 ON channels(peak_viewers);

CREATE TABLE IF NOT EXISTS categories(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    banner TEXT NOT NULL,
    live_viewers INTEGER NOT NULL,
    live_channels INTEGER NOT NULL,
    peak_viewers INTEGER NOT NULL,
    peak_channels INTEGER NOT NULL,
    description TEXT NOT NULL
);

CREATE INDEX categories_idx_1 ON categories(peak_viewers);
CREATE INDEX categories_idx_2 ON categories(peak_channels);

CREATE TABLE IF NOT EXISTS clips(
    id TEXT NOT NULL UNIQUE,
    category_name TEXT NOT NULL,
    category_slug TEXT NOT NULL,
    channel_name TEXT NOT NULL,
    channel_slug TEXT NOT NULL,
    is_mature BOOLEAN NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    likes INTEGER NOT NULL,
    livestream_id TEXT NOT NULL,
    thumbnail TEXT NOT NULL,
    views INTEGER NOT NULL,
    duration INTEGER NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE INDEX clips_idx_1 ON clips(likes);
CREATE INDEX clips_idx_2 ON clips(views);
CREATE INDEX clips_idx_3 ON clips(created_at);

CREATE TABLE IF NOT EXISTS overall_live_channels_chart(
    ts TEXT NOT NULL,
    n INTEGER NOT NULL
);

CREATE INDEX overall_live_channels_chart_idx_1 ON overall_live_channels_chart(ts, n);

CREATE TABLE IF NOT EXISTS overall_viewers_chart(
    ts TEXT NOT NULL,
    n INTEGER NOT NULL
);

CREATE INDEX overall_viewers_chart_idx_1 ON overall_viewers_chart(ts, n);

CREATE TABLE IF NOT EXISTS channels_followers_chart(
    slug TEXT NOT NULL,
    ts TEXT NOT NULL,
    n INTEGER NOT NULL
);

CREATE INDEX channels_followers_chart_idx_1 ON channels_followers_chart(ts, n);

CREATE TABLE IF NOT EXISTS channels_views_chart(
    slug TEXT NOT NULL,
    ts TEXT NOT NULL,
    n INTEGER NOT NULL
);

CREATE INDEX channels_views_chart_idx_1 ON channels_views_chart(ts, n);


CREATE TABLE IF NOT EXISTS categories_live_channels_chart(
    slug TEXT NOT NULL,
    ts TEXT NOT NULL,
    n INTEGER NOT NULL
);

CREATE INDEX categories_live_channels_chart_idx_1 ON categories_live_channels_chart(ts, n);

CREATE TABLE IF NOT EXISTS categories_views_chart(
    slug TEXT NOT NULL,
    ts TEXT NOT NULL,
    n INTEGER NOT NULL
);

CREATE INDEX categories_views_chart_idx_1 ON categories_views_chart(ts, n);