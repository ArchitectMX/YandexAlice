from flask import Flask, request, jsonify
import time

app = Flask(__name__)

user_states = {}

@app.route('/post', methods=['POST'])
def handle_alice():
    req = request.get_json()

    user_id = req['session']['user_id']
    if user_id not in user_states:
        user_states[user_id] = {}

    user = user_states[user_id]

    res = {
        "version": req['version'],
        "session": req['session'],
        "response": {
            "text": "",
            "buttons": []
        }
    }

    handle_dialog(req, res, user)
    return jsonify(res)


def handle_dialog(req, res, user):
    input_text = req['request']['original_utterance'].lower()

    if not input_text.strip():
        res['response']['text'] = "Привет! Чтобы начать, скажите 'да', 'согласен' или 'покажи PA'."
        return

    if any(keyword in input_text for keyword in ["покажи pa", "полный pa", "скажи pa", "полное соглашение", "police agreement", "договор", "условия"]):
        res['response']['text'] = """
        🚨 Police Agreement
        1. Неприкосновенность и неразглашение автора продукта
        Пользователь обязуется уважать право интеллектуальной собственности автора данного продукта.
        Запрещается:
        - Копирование, распространение, изменение или публикация кода/контента без согласия автора.
        - Передача информации о внутреннем устройстве продукта третьим лицам, включая, но не ограничиваясь: друзьям, коллегам, чатам в Telegram и бабушке.

        Нарушение данного пункта может расцениваться как акт цифрового хамства и будет встречено строго.
        2. Обязанности пользователя
        Пользователь, принимая настоящее соглашение, подтверждает, что:
        - Ознакомлен с функциональностью продукта.
        - Использует его в рамках правил, установленных автором.
        - Не делает глупостей, включая взлом, дизассемблирование или деструктивные действия.
        3. Ответственность за нарушение условий
        В случае невыполнения вышеуказанных обязанностей, пользователь:
        - Подлежит санкциям в виде одного (1) мега сырка от РосАгроКомплекса, передаваемого автору в натуральной или мемной форме.
        - Может быть также внесён в условный список "подозрительных личностей".
        4. Заключительные положения
        Принятие данного соглашения осуществляется автоматически при первом использовании продукта.
        Если вы не согласны с условиями — пожалуйста, закройте вкладку и не нарушайте наш внутренний дзен.
        """
        return

    # Этап 1: Согласие с Police Agreement
    if 'agreed_to_pa' not in user:
        if any(agreement in input_text for agreement in ["да", "согласен", "окей", "подтверждаю"]):
            user['agreed_to_pa'] = True
            res['response']['text'] = "Спасибо за согласие с Police Agreement. Перейдём к следующему этапу!"
            return

        res['response']['text'] = """Перед началом вы должны согласиться с Police Agreement. В нём:
        - Неприкосновенность автора
        - Обязанности по созданию навыка
        - В случае нарушения: 1 мега сырок от РосАгроКомплекса.
        Согласны?"""
        res['response']['buttons'] = [{"title": "Да", "hide": True}, {"title": "Согласен", "hide": True}]
        return

    # Отказ после согласия на создание навыка
    if 'agreed_to_make_skill' in user and user['agreed_to_make_skill']:
        if any(reject in input_text for reject in ["не хочу", "не буду", "отказ", "передумал"]):
            res['response']['text'] = "Согласно Police Agreement, вы обязаны выполнить создание навыка. В противном случае вас ждёт 1 мега сырок от РосАгроКомплекса."
            res['response']['buttons'] = [{"title": "Ладно, сделаю", "hide": True}, {"title": "Прошу прощения", "hide": True}]
            return

    # Согласие на создание навыка
    if "сделаю к следующему уроку" in input_text:
        user['agreed_to_make_skill'] = True
        days, hours, minutes, seconds = time_left_until_deadline()
        res['response']['text'] = f"Хорошо, жду. У вас есть: {format_time(days, hours, minutes, seconds)}"
        return

    # Запрос времени
    if any(time_query in input_text for time_query in ["сколько осталось времени", "сколько времени осталось", "сколько прошло", "сколько прошло времени"]):
        days, hours, minutes, seconds = time_left_until_deadline()
        res['response']['text'] = f"Прошло: 0 секунд. Осталось: {format_time(days, hours, minutes, seconds)}"
        return

    # Финальное согласие
    if any(accept in input_text for accept in ["ладно", "сделаю", "хорошо"]):
        res['response']['text'] = "Жду ваш навык!"
        res['response']['end_session'] = True
        return

    # Общий случай — агитация
    res['response']['text'] = f"Все говорят \"{req['request']['original_utterance']}\", а ты сделай навык для Алисы!"
    res['response']['buttons'] = [{"title": "Сделаю к следующему уроку", "hide": True}, {"title": "Ладно", "hide": True}]


def time_left_until_deadline():
    now = time.localtime()
    target = time.localtime(time.mktime(now) + (3 * 24 * 60 * 60))  # Примерная цель через 3 дня
    diff = time.mktime(target) - time.mktime(now)
    days = diff // (24 * 60 * 60)
    hours = (diff % (24 * 60 * 60)) // 3600
    minutes = (diff % 3600) // 60
    seconds = diff % 60
    return int(days), int(hours), int(minutes), int(seconds)


def format_time(days, hours, minutes, seconds):
    return f"{days} дн., {hours} ч., {minutes} мин., {seconds} сек."


if __name__ == '__main__':
    app.run(debug=True)
