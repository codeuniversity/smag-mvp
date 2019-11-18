import face_recognition
import tempfile
import os
import requests


def recognize(url):
    image = download_and_read_image(url)
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
    with open('img.jpg', 'wb') as handle:
        response = requests.get(url, stream=True)

        if not response.ok:
            print(response)
            raise Exception("coudln't download file: ", response)

        for block in response.iter_content(1024):
            if not block:
                break

            handle.write(block)
    image = face_recognition.load_image_file('img.jpg')
    os.remove('img.jpg')
    return image
