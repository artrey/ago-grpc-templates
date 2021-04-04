-- табличка с шаблонами
CREATE TABLE templates
(
    id      BIGSERIAL PRIMARY KEY,
    title   VARCHAR(100) NOT NULL,
    phone   VARCHAR(12)  NOT NULL,
    created TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- триггер для обновления поля updated в таблице шаблонов
CREATE TRIGGER update_timestamp
    BEFORE UPDATE
    ON templates
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_update_timestamp();
