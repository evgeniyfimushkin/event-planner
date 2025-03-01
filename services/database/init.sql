\c postgres;
DROP DATABASE IF EXISTS auth_db;
CREATE DATABASE auth_db OWNER postgres;

DROP DATABASE IF EXISTS users_db;
CREATE DATABASE users_db OWNER postgres;

DROP DATABASE IF EXISTS events_db;
CREATE DATABASE events_db OWNER postgres;

DROP DATABASE IF EXISTS registrations_db;
CREATE DATABASE registrations_db OWNER postgres;

DROP DATABASE IF EXISTS reviews_db;
CREATE DATABASE reviews_db OWNER postgres;

DROP DATABASE IF EXISTS media_db;
CREATE DATABASE media_db OWNER postgres;

DROP DATABASE IF EXISTS chat_db;
CREATE DATABASE chat_db OWNER postgres;

DROP DATABASE IF EXISTS notifications_db;
CREATE DATABASE notifications_db OWNER postgres;

