import face_recognition
import tempfile
import os
import requests
import random
import string


def recognize(url):
    image = download_and_read_image(url)
    if image is None:
        return []

    locations = face_recognition.face_locations(image)
    encodings = face_recognition.face_encodings(
        image, known_face_locations=locations)

    img_height = len(image)
    img_width = len(image[0])

    faces = []
    for index, location in enumerate(locations):
        encoding = encodings[index]
        (top, right, bottom, left) = location
        x = left
        y = top
        width = right - left
        height = bottom - top

        area = {
            "x": x,
            "y": y,
            "width": width,
            "height": height,
        }
        faces.append({
            "area": area,
            "encoding": encoding,
        })
    return faces


def download_and_read_image(url):
    file_name = random_string() + ".jpg"
    with open(file_name, 'wb') as handle:
        response = requests.get(url, stream=True)

        if not response.ok:
            print("failed to get file with name: " + file_name)
            return None

        for block in response.iter_content(1024):
            if not block:
                break

            handle.write(block)
    image = face_recognition.load_image_file(file_name)

    os.remove(file_name)

    return image


def random_string(n=12):
    return ''.join(random.choices(string.ascii_uppercase + string.digits, k=n))
