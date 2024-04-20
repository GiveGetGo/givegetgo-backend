import React, { useState } from 'react';
import { View, StyleSheet, Image, TouchableOpacity } from 'react-native';
import { TextInput, Button, Appbar } from 'react-native-paper';

const CreatePostScreen: React.FC = () => {
  const [title, setTitle] = useState<string>('');
  const [description, setDescription] = useState<string>('');

  const handleImageUpload = () => {
    // TODO: Implement your image upload logic here
  };

  const handleSubmitPost = () => {
    // TODO: Implement your submit post logic here
    console.log('Post submitted with Title:', title, 'and Description:', description);
  };

  // Placeholder for your image icon, replace 'require' with the actual path to your icon
  const imageUploadIcon = require('./image_icon.jpg');

  return (
    <View style={styles.container}>
      <Appbar.Header>
        <Appbar.Content title="GiveGetGo" />
        <Appbar.Action icon="magnify" onPress={() => {}} />
      </Appbar.Header>
      
      <TextInput
        label="Add a title..."
        value={title}
        onChangeText={setTitle}
        style={styles.input}
      />

      <TextInput
        label="Add a description..."
        value={description}
        onChangeText={setDescription}
        style={styles.input}
        multiline={true}
      />

      <TouchableOpacity onPress={handleImageUpload} style={styles.imageUploadButton}>
        <Image source={imageUploadIcon} style={styles.imageIcon} />
      </TouchableOpacity>

      <Button mode="contained" onPress={handleSubmitPost} style={styles.button}>
        Post
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 8,
  },
  input: {
    marginVertical: 8,
  },
  imageUploadButton: {
    // Style for your image upload button
    alignItems: 'center',
    justifyContent: 'center',
    padding: 16,
    marginVertical: 8,
  },
  imageIcon: {
    width: 50, // Adjust according to your icon size
    height: 50, // Adjust according to your icon size
  },
  button: {
    marginVertical: 8,
  },
  // Add other styles as needed
});

export default CreatePostScreen;
